package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/wjoseperez20/boletia-currency-api/pkg/cache"
	"github.com/wjoseperez20/boletia-currency-api/pkg/database"
	"github.com/wjoseperez20/boletia-currency-api/pkg/models"
)

var httpClient *http.Client

func init() {
	// Initialize HTTP client with timeout and transport settings
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

// getCurrencyData retrieves currency data from an external API.
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

	// Measure request time
	start := time.Now()

	resp, err := httpClient.Do(req)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			// Log request history
			requestLog := models.RequestHistory{
				Endpoint:     apiEndpoint,
				ResponseTime: 0,
				StatusCode:   http.StatusRequestTimeout,
			}

			if err := insertRequestHistory(requestLog); err != nil {
				log.Printf("Error inserting request history: %s\n", err)
			}
		}

		return models.CurrencyAPIResponse{}, fmt.Errorf("error sending HTTP request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %s\n", err)
		}
	}(resp.Body)
	elapsed := time.Since(start).Seconds()

	// Log request history
	requestLog := models.RequestHistory{
		Endpoint:     apiEndpoint,
		ResponseTime: elapsed,
		StatusCode:   resp.StatusCode,
	}
	if err := insertRequestHistory(requestLog); err != nil {
		log.Printf("Error inserting request history: %s\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		return models.CurrencyAPIResponse{}, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var currencyResponse models.CurrencyAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&currencyResponse); err != nil {
		return models.CurrencyAPIResponse{}, fmt.Errorf("error decoding JSON response: %v", err)
	}

	return currencyResponse, nil
}

// insertCurrencies inserts currency data into the database.
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
			Name:      data.Code,
			Code:      code,
			Value:     data.Value,
			CreatedAt: currencyResponse.Meta.LastUpdatedAt,
		}

		if err := tx.Create(&currency).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error inserting currency data: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing database transaction: %v", err)
	}

	if err := invalidateCache(); err != nil {
		return fmt.Errorf("error invalidating cache: %v", err)
	}

	return nil
}

// invalidateCache deletes cached currency data.
func invalidateCache() error {
	keysPattern := "currency_*"
	keys, err := cache.Rdb.Keys(cache.Ctx, keysPattern).Result()
	if err != nil {
		return fmt.Errorf("error retrieving cache keys: %v", err)
	}

	for _, key := range keys {
		if err := cache.Rdb.Del(cache.Ctx, key).Err(); err != nil {
			return fmt.Errorf("error deleting cache key %s: %v", key, err)
		}
	}

	return nil
}

func insertRequestHistory(requestLog models.RequestHistory) error {
	db := database.DB

	if err := db.Create(&requestLog).Error; err != nil {
		return fmt.Errorf("error inserting request history: %v", err)
	}

	return nil
}

// InitDaemon initializes the daemon for periodic currency updates.
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

	// Wait indefinitely to prevent the function from exiting.
	select {}
}
