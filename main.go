package main

import (
	"log"
	"url-shortener/domains/urls/routes"
	infrastructure "url-shortener/infrastructures"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := infrastructure.OpenPostgreConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	router := gin.Default()
	routes.SetUrlShortenRoutes(db, router)
	router.Run(":8080")
}
