package models

import "gorm.io/gorm"

type Shortlink struct {
	gorm.Model
	ID       uint64 `json:"id" gorm:"primaryKey"`
	UserID   uint32 `json:"user_id"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}
