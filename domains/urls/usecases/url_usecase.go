package usecases

import (
	"crypto/rand"
	"log"
	"time"
	"url-shortener/domains/urls/entities"
	"url-shortener/domains/urls/interfaces"
	"url-shortener/shared"
)

type UrlMappingUsecase struct {
	UrlMappingRepo interfaces.UrlMappingInterface
	UrlClickRepo   interfaces.UrlClickInterface
}

func NewUrlMappingUsecase(urlMappingRepo interfaces.UrlMappingInterface, urlClickRepo interfaces.UrlClickInterface) *UrlMappingUsecase {
	return &UrlMappingUsecase{
		UrlMappingRepo: urlMappingRepo,
		UrlClickRepo:   urlClickRepo,
	}
}

func (uc *UrlMappingUsecase) ShortenUrl(longUrl string, expiresAt *time.Time) (*entities.UrlMappingEntity, error) {
	if expiresAt == nil {
		defaultExpiration := time.Now().Add(5 * time.Hour)
		expiresAt = &defaultExpiration
	}

	existing, err := uc.UrlMappingRepo.GetByLongUrl(longUrl)
	if err != nil {
		return nil, shared.NewBadRequestError("URL already exists")
	}

	if existing != nil {
		if existing.ExpiresAt != nil && existing.ExpiresAt.Before(time.Now()) {
			// If URL is expired, generate new short code and create new mapping
			shortCode, err := uc.generateRandomCode(6)
			if err != nil {
				return nil, shared.NewInternalServerError("Failed to generate short code: " + err.Error())
			}

			mapping, err := uc.UrlMappingRepo.CreateShortUrl(longUrl, shortCode, expiresAt)
			if err != nil {
				return nil, shared.NewInternalServerError("Failed to create short URL: " + err.Error())
			}

			return mapping, nil
		} else {
			return nil, shared.NewBadRequestError("URL already exists and not expired")
		}
	}

	// Generate a random short code
	shortCode, err := uc.generateRandomCode(6)
	if err != nil {
		return nil, shared.NewInternalServerError("Failed to generate short code: " + err.Error())
	}

	// Create new mapping
	mapping, err := uc.UrlMappingRepo.CreateShortUrl(longUrl, shortCode, expiresAt)
	if err != nil {
		return nil, shared.NewInternalServerError("Failed to create short URL: " + err.Error())
	}

	return mapping, nil
}

func (uc *UrlMappingUsecase) GetOriUrlByShortCode(shortCode string) (*entities.UrlMappingEntity, error) {
	if shortCode == "" {
		return nil, shared.NewBadRequestError("Short code cannot be empty")
	}

	mapping, err := uc.UrlMappingRepo.GetByShortCode(shortCode)
	if err != nil {
		if err.Error() == "short URL not found" {
			return nil, shared.NewNotFoundError("Short URL not found")
		}
		if err.Error() == "short URL has expired" {
			return nil, shared.NewBadRequestError("Short URL has expired")
		}
		return nil, shared.NewInternalServerError("Failed to retrieve short URL: " + err.Error())
	}

	return mapping, nil
}

func (uc *UrlMappingUsecase) GetOriUrlByShortCode2(shortCode string, ipAddress string, userAgent string) (*entities.UrlMappingEntity, error) {
	if shortCode == "" {
		return nil, shared.NewBadRequestError("Short code cannot be empty")
	}

	mapping, err := uc.UrlMappingRepo.GetByShortCode(shortCode)
	if err != nil {
		if err.Error() == "short URL not found" {
			return nil, shared.NewNotFoundError("Short URL not found")
		}
		if err.Error() == "short URL has expired" {
			return nil, shared.NewBadRequestError("Short URL has expired")
		}
		return nil, shared.NewInternalServerError("Failed to retrieve short URL: " + err.Error())
	}

	// Log click
	_, err = uc.UrlClickRepo.LogClick(mapping.ID, ipAddress, userAgent)
	if err != nil {
		log.Printf("Error logging click - ShortCode: %s, MappingID: %d, IPAddress: %s, Error: %v",
			shortCode, mapping.ID, ipAddress, err)
	} else {
		log.Printf("Successfully logged click - ShortCode: %s, MappingID: %d, IPAddress: %s",
			shortCode, mapping.ID, ipAddress)
	}

	return mapping, nil
}

func (uc *UrlMappingUsecase) generateRandomCode(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes), nil
}
