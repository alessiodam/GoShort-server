package database

import (
	"GoShort/internal/models"
	"GoShort/logging"
	"errors"
	"log"
	"time"
)

var (
	sessionLogger = logging.NewLogger("database.session")
)

func CreateSession(session *models.Session) error {
	if err := instance.Create(session).Error; err != nil {
		sessionLogger.Warningf("Failed to create session: %v", err)
		return err
	}
	return nil
}

func GetUserBySessionToken(token string) (*models.User, error) {
	var session models.Session
	if err := instance.Where("token = ?", token).First(&session).Error; err != nil {
		sessionLogger.Warningf("Failed to get session: %v", err)
		return nil, errors.New("session not found")
	}
	if session.Expired {
		sessionLogger.Warningf("Session expired: %v", session)
		return nil, errors.New("session has expired")
	}
	if session.ExpiresAt.Before(time.Now()) {
		sessionLogger.Warningf("Session expired: %v", session)
		session.Expired = true
		if err := instance.Save(&session).Error; err != nil {
			return nil, err
		}
		return nil, errors.New("session has expired")
	}

	user, err := GetUserByID(session.UserID)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return nil, errors.New("user not found")
	}
	return user, nil
}

func DeleteSessionByToken(token string) error {
	if err := instance.Where("token = ?", token).Delete(&models.Session{}).Error; err != nil {
		sessionLogger.Warningf("Failed to delete session: %v", err)
		return err
	}
	return nil
}
