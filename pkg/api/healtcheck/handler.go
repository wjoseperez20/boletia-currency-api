package healtcheck

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// Healthcheck godoc
// @Summary Healthcheck
// @Schemes
// @Description do ping
// @Tags Healthcheck
// @Accept json
// @Produce json
// @Success 200 {string} ok
// @Router /_ [get]
func Healthcheck(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"message": "ok"})
}
