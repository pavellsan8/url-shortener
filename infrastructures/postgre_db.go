package infrastructure

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"url-shortener/config"
)

func OpenPostgreConnection() (*gorm.DB, error) {
	cfg := config.GetConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Database.Host, cfg.Database.Username, cfg.Database.Password, cfg.Database.Name,
		cfg.Database.Port, cfg.Database.SslMode, cfg.Database.TimeZone)

	fmt.Println("Connecting to database with DSN:", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	fmt.Println("Database connection established successfully")
	return db, nil
}
