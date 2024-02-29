package api

import (
	"github.com/wjoseperez20/boletia-currency-api/docs"
	"github.com/wjoseperez20/boletia-currency-api/pkg/api/currencies"
	"github.com/wjoseperez20/boletia-currency-api/pkg/api/healtcheck"
	"github.com/wjoseperez20/boletia-currency-api/pkg/api/users"
	"github.com/wjoseperez20/boletia-currency-api/pkg/middleware"
	"time"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.Security())
		r.Use(middleware.Xss())
	}
	r.Use(middleware.Cors())
	r.Use(middleware.RateLimiter(rate.Every(1*time.Minute), 60)) // 60 requests per minute

	// api routes
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		v1.GET("/_", healtcheck.Healthcheck)
		v1.POST("/login", middleware.APIKeyAuth(), users.LoginUser)
		v1.POST("/register", middleware.APIKeyAuth(), users.RegisterUser)

		// Currency
		v1.GET("/currencies", middleware.JWTAuth(), currencies.FindCurrencies)
		v1.GET("/currencies/:name", middleware.JWTAuth(), currencies.FindCurrency)
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
