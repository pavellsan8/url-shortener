package entities

import "time"

type UrlMappingEntity struct {
	ID        uint      `gorm:"primaryKey"`
	ShortCode string    `gorm:"column:short_code;size:16;not null;uniqueIndex"`
	LongURL   string    `gorm:"column:long_url;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ExpiresAt *time.Time
}

func (UrlMappingEntity) TableName() string {
	return "url_mappings"
}
