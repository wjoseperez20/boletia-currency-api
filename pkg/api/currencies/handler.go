package currencies

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wjoseperez20/boletia-currency-api/pkg/cache"
	"github.com/wjoseperez20/boletia-currency-api/pkg/database"
	"github.com/wjoseperez20/boletia-currency-api/pkg/models"
)

func Manager(c *gin.Context) {
	// Get query params
	currencyName := strings.ToUpper(c.Param("name"))
	finitQuery := c.DefaultQuery("finit", "")
	fendQuery := c.DefaultQuery("fend", "")

	// Check if currency name is empty
	if currencyName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid currency param"})
		return
	}

	// Take action based on query params
	if currencyName == "ALL" {
		FetchAllCurrencies(c)
		return
	}

	// Parse query params into time.Time
	var finit, fend time.Time
	var err error
	layout := "2006-01-02T15:04:05"

	if finitQuery != "" {
		finit, err = time.Parse(layout, finitQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid param start date format"})
			return
		}
	}

	if fendQuery != "" {
		fend, err = time.Parse(layout, fendQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid param end date format"})
			return
		}
	}

	// Fetch or retrieve currencies
	FindCurrencyByDateRange(c, currencyName, finit, fend, finitQuery, fendQuery)
}

func FetchAllCurrencies(c *gin.Context) {
	// Get all currencies from the database or cache
	var currencies []models.Currency
	cacheKey := "currency_all"

	// Attempt to retrieve currencies from cache
	if cachedCurrencies, err := getCurrenciesFromCache(cacheKey); err == nil {
		c.JSON(http.StatusOK, cachedCurrencies)
		return
	}

	// Fetch currencies from the database
	database.DB.Find(&currencies)

	// Store currencies in cache
	storeCurrenciesInCache(cacheKey, currencies)

	c.JSON(http.StatusOK, currencies)
}

func FindCurrencyByDateRange(c *gin.Context, currencyName string, finit, fend time.Time, finitQuery, fendQuery string) {
	// Get currencies by date range from the database or cache
	var currencies []models.Currency
	cacheKey := "currency_" + currencyName + "_finit_" + finitQuery + "_fend_" + fendQuery

	// Attempt to retrieve currencies from cache
	if cachedCurrencies, err := getCurrenciesFromCache(cacheKey); err == nil {
		c.JSON(http.StatusOK, cachedCurrencies)
		return
	}

	// Construct database query for currency and date range
	dbQuery := database.DB.Where("name = ?", currencyName)

	if !finit.IsZero() && !fend.IsZero() {
		dbQuery = dbQuery.Where("created_at BETWEEN ? AND ?", finit, fend)
	} else if !finit.IsZero() {
		dbQuery = dbQuery.Where("created_at >= ?", finit)
	} else if !fend.IsZero() {
		dbQuery = dbQuery.Where("created_at <= ?", fend)
	}

	// Fetch currencies from the database
	dbQuery.Find(&currencies)

	// Store currencies in cache
	storeCurrenciesInCache(cacheKey, currencies)

	c.JSON(http.StatusOK, currencies)
}

func getCurrenciesFromCache(cacheKey string) ([]models.Currency, error) {
	// Retrieve currencies from cache
	var currencies []models.Currency

	cachedCurrencies, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
	if err != nil {
		return currencies, err
	}

	err = json.Unmarshal([]byte(cachedCurrencies), &currencies)
	return currencies, err
}

func storeCurrenciesInCache(cacheKey string, currencies []models.Currency) {
	// Store currencies in cache
	serializedCurrencies, err := json.Marshal(currencies)
	if err != nil {
		return
	}

	_ = cache.Rdb.Set(cache.Ctx, cacheKey, serializedCurrencies, 0).Err()
}
