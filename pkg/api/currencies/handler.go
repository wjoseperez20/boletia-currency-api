package currencies

import (
	"github.com/gin-gonic/gin"
	"github.com/wjoseperez20/boletia-currency-api/pkg/database"
	"github.com/wjoseperez20/boletia-currency-api/pkg/models"
	"net/http"
	"strconv"
)

// @BasePath /api/v1
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

	if err := database.DB.Where("currency = ?", c.Param("currency")).First(&currency).Error; err != nil {
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

	// Get query params
	offsetQuery := c.DefaultQuery("offset", "0")
	limitQuery := c.DefaultQuery("limit", "10")

	// Convert query params to integers
	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset format"})
		return
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit format"})
		return
	}

	// Get currencies with pagination
	if err := database.DB.Offset(offset).Limit(limit).Find(&currencies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching currencies"})
		return
	}

	c.JSON(http.StatusOK, currencies)

}
