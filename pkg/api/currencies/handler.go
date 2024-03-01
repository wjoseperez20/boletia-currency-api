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

// HandleCurrencyRequest godoc
// @Summary Manage currency requests
// @Description check param name to get all currencies or a specific currency by date range
// @Produce json
// @Param name path string true "Currency name"
func HandleCurrencyRequest(c *gin.Context) {
	// Get query params
	currencyName := strings.ToUpper(c.Param("name"))
	finitQuery := c.DefaultQuery("finit", "")
	fendQuery := c.DefaultQuery("fend", "")

	// Check if currency name is "ALL" to fetch all currencies
	// The Albanian Lek (ALL) code conflicts with the "ALL" keyword
	if currencyName == "ALL" && finitQuery == "" && fendQuery == "" {
		fetchAllCurrencies(c)
		return
	}

	// Check database for the currency
	var currency models.Currency
	if err := database.DB.Where("name = ?", currencyName).First(&currency).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Currency is not valid"})
		return
	}

	// Parse query params into time.Time
	layout := "2006-01-02T15:04:05"
	var finit, fend time.Time
	var err error

	if finitQuery != "" {
		finit, err = time.Parse(layout, finitQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid finit date format"})
			return
		}
	}

	if fendQuery != "" {
		fend, err = time.Parse(layout, fendQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fend date format"})
			return
		}
	}

	// Fetch or retrieve currencies by date range
	fetchCurrencyByDateRange(c, currencyName, finit, fend)
}

// fetchAllCurrencies godoc
// @Summary Get all currencies
// @Description Get all currencies from the database
// @Tags Currencies
// @Produce json
// @Success 200 {object} []models.GroupedCurrencies
// @Failure 404 {string} string "No currencies found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /currencies/all [get]
func fetchAllCurrencies(c *gin.Context) {
	// Get all currencies from the database or cache
	var groupedCurrencies []models.GroupedCurrencies
	cacheKey := "currency_all"

	// Attempt to retrieve currencies from cache
	if cachedCurrencies, err := getCurrenciesFromCache(cacheKey); err == nil {
		c.JSON(http.StatusOK, cachedCurrencies)
		return
	}

	// Get all currencies from the database
	var currencies []models.Currency
	if err := database.DB.Select("name, created_at, value").Find(&currencies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(currencies) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No currencies found"})
		return

	}

	// Create a map to store grouped currencies
	currencyMap := make(map[string][]models.Currency)

	// Group currencies by code
	for _, currency := range currencies {
		currencyMap[currency.Name] = append(currencyMap[currency.Name], currency)
	}

	// Format grouped currencies into GroupedCurrencies struct
	for code, currencyList := range currencyMap {
		var data []models.CurrencyData

		for _, currency := range currencyList {
			data = append(data, models.CurrencyData{
				Date:  currency.CreatedAt.Format("2006-01-02T15:04:05"),
				Value: currency.Value,
			})
		}

		groupedCurrencies = append(groupedCurrencies, models.GroupedCurrencies{
			Code: code,
			Data: data,
		})
	}

	// Store currencies in cache
	storeCurrenciesInCache(cacheKey, groupedCurrencies)

	c.JSON(http.StatusOK, groupedCurrencies)
}

// fetchCurrencyByDateRange godoc
// @Summary Get currency by date range
// @Description Get a specific currency by date range from the database
// @Produce json
// @Tags Currencies
// @Param name path string true "Currency name"
// @Param finit query string false "Start date"
// @Param fend query string false "End date"
// @Success 200 {object} models.GroupedCurrencies
// @Failure 404 {string} string "No currencies found for the specified date range"
// @Failure 500 {string} string "Internal Server Error"
// @Router /currencies/{name} [get]
func fetchCurrencyByDateRange(c *gin.Context, currencyName string, startDate, endDate time.Time) {
	// Prepare cache key using currency name and date range
	cacheKey := "currency_" + currencyName + "_start_" + startDate.Format("2006-01-02T15:04:05") + "_end_" + endDate.Format("2006-01-02T15:04:05")

	// Attempt to retrieve currencies from cache
	if cachedCurrencies, err := getCurrenciesFromCache(cacheKey); err == nil {
		c.JSON(http.StatusOK, cachedCurrencies)
		return
	}

	// Retrieve currency history from the database
	var currencyHistory []models.Currency
	if err := database.DB.Select("name, created_at, value").
		Where("name = ? AND created_at BETWEEN ? AND ?", currencyName, startDate, endDate).
		Find(&currencyHistory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(currencyHistory) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No currencies found for the specified date range"})
		return

	}

	// Format currency history into desired structure
	var formattedHistory []models.CurrencyData
	for _, history := range currencyHistory {
		formattedHistory = append(formattedHistory, models.CurrencyData{
			Date:  history.CreatedAt.Format("2006-01-02T15:04:05"),
			Value: history.Value,
		})
	}
	groupedCurrencies := models.GroupedCurrencies{
		Code: currencyName,
		Data: formattedHistory,
	}

	// Store currency history in cache
	storeCurrenciesInCache(cacheKey, groupedCurrencies)

	c.JSON(http.StatusOK, groupedCurrencies)
}

func getCurrenciesFromCache(cacheKey string) (interface{}, error) {
	// Retrieve currencies from cache
	var currencies interface{}

	cachedCurrencies, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
	if err != nil {
		return currencies, err
	}

	err = json.Unmarshal([]byte(cachedCurrencies), &currencies)
	return currencies, err
}

func storeCurrenciesInCache(cacheKey string, currency interface{}) {
	// Store currencies in cache
	serializedCurrencies, err := json.Marshal(currency)
	if err != nil {
		return
	}

	_ = cache.Rdb.Set(cache.Ctx, cacheKey, serializedCurrencies, 0).Err()
}
