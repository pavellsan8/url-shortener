package responses

import "time"

type UrlMappingResponse struct {
	ShortCode string     `json:"short_code"`
	LongUrl   string     `json:"long_url"`
	ExpiresAt *time.Time `json:"expires_at"`
	FullUrl   string     `json:"short_url"`
}
