package repositories

import (
	"log"
	"time"
	"url-shortener/domains/urls/entities"

	"gorm.io/gorm"
)

type UrlClickRepository struct {
	db *gorm.DB
}

func NewUrlClickRepository(db *gorm.DB) *UrlClickRepository {
	return &UrlClickRepository{db: db}
}

func (r *UrlClickRepository) LogClick(mappingID uint, ipAddress, userAgent string) (*entities.UrlClickEntity, error) {
	click := &entities.UrlClickEntity{
		MappingID: mappingID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ClickedAt: time.Now(),
	}

	if err := r.db.Create(&click).Error; err != nil {
		log.Printf("Error logging click: %v", err)
		return nil, err
	}

	return click, nil
}
