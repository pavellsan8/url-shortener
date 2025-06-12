package routes

import (
	"url-shortener/domains/urls/handlers"
	"url-shortener/domains/urls/repositories"
	"url-shortener/domains/urls/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUrlShortenRoutes(db *gorm.DB, router *gin.Engine) {
	urlMappingRepo := repositories.NewUrlMappingRepository(db)
	urlClickRepo := repositories.NewUrlClickRepository(db)
	urlMappingUseCase := usecases.NewUrlMappingUsecase(urlMappingRepo, urlClickRepo)
	urlMappingHandler := handlers.NewUrlMappingHandler(urlMappingUseCase)

	api := router.Group("/api/v1")
	{
		api.POST("/shorten-url", func(c *gin.Context) { urlMappingHandler.ShortenUrl(c.Writer, c.Request) })
		api.GET("/get-long-url-data", func(c *gin.Context) { urlMappingHandler.GetLongUrlFromCode(c) })
	}

	router.GET("/:shortCode", func(c *gin.Context) {
		shortCode := c.Param("shortCode")
		if shortCode == "api" || shortCode == "favicon.ico" {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}
		urlMappingHandler.RedirectOriUrl(c.Writer, c.Request)
	})
}
