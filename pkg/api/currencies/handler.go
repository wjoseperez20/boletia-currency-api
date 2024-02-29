package currencies

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wjoseperez20/boletia-currency-api/pkg/cache"
	"github.com/wjoseperez20/boletia-currency-api/pkg/database"
	"github.com/wjoseperez20/boletia-currency-api/pkg/models"
	"net/http"
	"time"
)

// FindCurrency godoc
// @Summary Find a currency by ID
// @Description Get details of a currency by its ID
// @Tags Currencies
// @Security JwtAuth
// @Produce json
// @Param currency path string true "Currency ID"
// @Success 200 {object} models.Currency "Successfully retrieved currency"
// @Failure 404 {string} string "Currency not found"
// @Router /currencies/{currency} [get]
func FindCurrency(c *gin.Context) {
	var currency models.Currency

	// Check if currency exists
	if err := database.DB.Where("name = ?", c.Param("name")).First(&currency).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "currency not found"})
		return
	}

	c.JSON(http.StatusOK, currency)
}

// FindCurrencies godoc
// @Summary Get all currencies with pagination
// @Description Get a list of all currencies with optional pagination
// @Tags Currencies
// @Security JwtAuth
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {array} models.Currency "Successfully retrieved list of currencies"
// @Router /currencies [get]
func FindCurrencies(c *gin.Context) {
	var currencies []models.Currency

	// Get current date and format it to "YYYY-MM-DDT00:00:00Z"
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayFormatted := today.Format(time.RFC3339)

	// Get query params
	finitQuery := c.DefaultQuery("finit", todayFormatted)
	fendQuery := c.DefaultQuery("fend", todayFormatted)

	// Parse query params into time.Time
	finit, err := time.Parse(time.RFC3339, finitQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	fend, err := time.Parse(time.RFC3339, fendQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	// Create a cache key based on query params
	cacheKey := "currencies_finit_" + finitQuery + "_fend_" + fendQuery

	// Get currencies from cache
	cachedCurrencies, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
	if err == nil {
		err := json.Unmarshal([]byte(cachedCurrencies), &currencies)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching currencies"})
			return
		}
		c.JSON(http.StatusOK, currencies)
		return
	}

	// If cache missed, fetch data from the database
	database.DB.Where("created_at BETWEEN ? AND ?", finit, fend).Find(&currencies)

	// Serialize currencies object and store in cache
	serializedCurrencies, err := json.Marshal(currencies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching currencies"})
		return
	}
	err = cache.Rdb.Set(cache.Ctx, cacheKey, serializedCurrencies, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching currencies"})
		return
	}

	c.JSON(http.StatusOK, currencies)
}
