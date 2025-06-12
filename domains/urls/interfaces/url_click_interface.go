package interfaces

import "url-shortener/domains/urls/entities"

type UrlClickInterface interface {
	LogClick(mappingID uint, ipAddress string, userAgent string) (*entities.UrlClickEntity, error)
}
