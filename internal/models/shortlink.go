package models

import "gorm.io/gorm"

type Shortlink struct {
	gorm.Model
	ID        uint32               `json:"id" gorm:"primaryKey"`
	UserID    uint32               `json:"user_id"`
	ShortURL  string               `json:"short_url"`
	LongURL   string               `json:"long_url"`
	Analytics []ShortlinkAnalytics `json:"analytics"`
}

type ShortlinkAnalytics struct {
	gorm.Model
	ShortlinkID uint32 `json:"shortlink_id"`
	Clicks      uint32 `json:"clicks"`
}

type ShortlinkBrowserAnalytics struct {
	gorm.Model
	ShortlinkID uint32 `json:"shortlink_id" gorm:"index"`
	Browser     string `json:"browser"`
	Country     string `json:"country"`
	Count       uint32 `json:"count"`
}
