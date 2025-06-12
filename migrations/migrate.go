package migrations

import (
	"fmt"
	"log"
	"url-shortener/config"
	"url-shortener/domains/urls/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func RunMigration() {
	cfg := config.GetConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Database.Host, cfg.Database.Username, cfg.Database.Password, cfg.Database.Name,
		cfg.Database.Port, cfg.Database.SslMode, cfg.Database.TimeZone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&entities.UrlMappingEntity{},
		&entities.UrlClickEntity{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed successfully!")
}
