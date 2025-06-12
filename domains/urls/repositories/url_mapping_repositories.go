package repositories

import (
	"errors"
	"time"
	"url-shortener/domains/urls/entities"

	"gorm.io/gorm"
)

type UrlMappingRepository struct {
	db *gorm.DB
}

func NewUrlMappingRepository(db *gorm.DB) *UrlMappingRepository {
	return &UrlMappingRepository{db: db}
}

func (r *UrlMappingRepository) CreateShortUrl(longUrl string, shortCode string, expiresAt *time.Time) (*entities.UrlMappingEntity, error) {
	mapping := &entities.UrlMappingEntity{
		ShortCode: shortCode,
		LongURL:   longUrl,
		ExpiresAt: expiresAt,
	}

	if err := r.db.Create(mapping).Error; err != nil {
		return nil, err
	}

	return mapping, nil
}

func (r *UrlMappingRepository) GetByShortCode(shortCode string) (*entities.UrlMappingEntity, error) {
	var mapping entities.UrlMappingEntity
	err := r.db.Where("short_code = ?", shortCode).First(&mapping).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("short URL not found")
		}
		return nil, err
	}

	if mapping.ExpiresAt != nil && mapping.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("short URL has expired")
	}

	return &mapping, nil
}

func (r *UrlMappingRepository) GetByLongUrl(longUrl string) (*entities.UrlMappingEntity, error) {
	var mapping entities.UrlMappingEntity
	err := r.db.Where("long_url = ?", longUrl).First(&mapping).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &mapping, nil
}

func (r *UrlMappingRepository) GetValidByShortCode(shortCode string) (*entities.UrlMappingEntity, error) {
	var mapping entities.UrlMappingEntity
	query := r.db.Where("short_code = ?", shortCode)

	query = query.Where("expires_at IS NULL OR expires_at > ?", time.Now())

	err := query.First(&mapping).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("short URL not found or expired")
		}
		return nil, err
	}

	return &mapping, nil
}
