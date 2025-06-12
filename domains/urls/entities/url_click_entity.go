package entities

import "time"

type UrlClickEntity struct {
	ID        uint      `gorm:"primaryKey"`
	MappingID uint      `gorm:"not null;index"`
	ClickedAt time.Time `gorm:"autoCreateTime"`
	IPAddress string    `gorm:"size:45"`
	UserAgent string    `gorm:"type:text"`
}

func (UrlClickEntity) TableName() string {
	return "url_clicks"
}
