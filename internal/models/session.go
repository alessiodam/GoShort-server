package models

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model

	ID        uint32    `json:"id" gorm:"primaryKey"`
	UserID    uint32    `json:"user_id"`
	Token     string    `json:"token" gorm:"unique"`
	Expired   bool      `json:"expired"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (session *Session) Expire() {
	session.Expired = true
	session.ExpiresAt = time.Now()
}

func (session *Session) IsExpired() bool {
	return session.ExpiresAt.Before(time.Now()) || session.Expired
}
