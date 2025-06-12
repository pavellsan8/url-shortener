package interfaces

import (
	"time"
	"url-shortener/domains/urls/entities"
)

type UrlMappingInterface interface {
	CreateShortUrl(shortCode string, longUrl string, expiresAt *time.Time) (*entities.UrlMappingEntity, error)
	GetByShortCode(shortCode string) (*entities.UrlMappingEntity, error)
	GetByLongUrl(longUrl string) (*entities.UrlMappingEntity, error)
	GetValidByShortCode(shortCode string) (*entities.UrlMappingEntity, error)
}
