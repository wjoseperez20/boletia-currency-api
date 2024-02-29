package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wjoseperez20/boletia-currency-api/pkg/database"
	"github.com/wjoseperez20/boletia-currency-api/pkg/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var httpClient *http.Client

func init() {
	timeoutStr := os.Getenv("CURRENCY_API_TIMEOUT")
	if timeoutStr == "" {
		log.Fatal("CURRENCY_API_TIMEOUT is not set")
	}

	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Fatalf("Error parsing CURRENCY_API_TIMEOUT: %v", err)
	}

	httpClient = &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConnsPerHost: 10,
		},
	}
}

func getCurrencyData() (models.CurrencyAPIResponse, error) {
	apiEndpoint := os.Getenv("CURRENCY_API_ENDPOINT")
	apiKey := os.Getenv("CURRENCY_API_KEY")

	if apiEndpoint == "" || apiKey == "" {
		return models.CurrencyAPIResponse{}, errors.New("missing API endpoint or API key")
	}

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return models.CurrencyAPIResponse{}, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Add("apikey", apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return models.CurrencyAPIResponse{}, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.CurrencyAPIResponse{}, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var currencyResponse models.CurrencyAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&currencyResponse); err != nil {
		return models.CurrencyAPIResponse{}, fmt.Errorf("error decoding JSON response: %v", err)
	}

	return currencyResponse, nil
}

func insertCurrencies(currencyResponse models.CurrencyAPIResponse) error {
	db := database.DB

	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("error starting database transaction: %v", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to panic")
		}
	}()

	for code, data := range currencyResponse.Data {
		currency := models.Currency{
			Name:  data.Code,
			Code:  code,
			Value: data.Value,
		}

		if err := tx.Create(&currency).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error inserting currency data: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing database transaction: %v", err)
	}

	return nil
}

func InitDaemon() {
	wakeupStr := os.Getenv("DAEMON_WAKEUP")
	if wakeupStr == "" {
		log.Fatal("DAEMON_WAKEUP is not set")
	}

	wakeup, err := strconv.Atoi(wakeupStr)
	if err != nil {
		log.Printf("Error parsing DAEMON_WAKEUP: %s\n", err)
		return
	}

	ticker := time.NewTicker(time.Duration(wakeup) * time.Second)

	go func() {
		for range ticker.C {
			currencyResponse, err := getCurrencyData()
			if err != nil {
				log.Printf("Error getting currency data: %s\n", err)
				continue
			}

			if err := insertCurrencies(currencyResponse); err != nil {
				log.Printf("Error inserting currency data: %s\n", err)
				continue
			}
		}
	}()
}
