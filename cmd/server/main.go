package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wjoseperez20/boletia-currency-api/pkg/api"
	"github.com/wjoseperez20/boletia-currency-api/pkg/cache"
	"github.com/wjoseperez20/boletia-currency-api/pkg/daemon"
	"github.com/wjoseperez20/boletia-currency-api/pkg/database"
	"log"
)

// @title           Boletia Currency API
// @version         1.0
// @description     This is a simple API for currencies exchange

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8001
// @BasePath  /api/v1

// @securityDefinitions.apikey JwtAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cache.InitRedis()
	database.ConnectDatabase()

	go daemon.InitDaemon()

	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	// Initialize Gin router
	r := api.InitRouter()

	// Run the server
	if err := r.Run(":8001"); err != nil {
		log.Fatal(err)
	}
}
