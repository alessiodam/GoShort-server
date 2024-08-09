package database

import (
	"GoShort/internal/models"
	"GoShort/logging"
	"errors"
	"gorm.io/gorm"
)

var (
	userLogger = logging.NewLogger("database.user")
)

func InsertUser(user *models.User) error {
	tx := instance.Begin()
	if tx.Error != nil {
		userLogger.Warningf("failed to begin transaction: %v", tx.Error)
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			userLogger.Errorf("recovered from panic in InsertUser: %v", r)
		}
	}()

	var existingUser models.User
	if err := tx.Where("email = ? OR username = ?", user.Email, user.Username).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return errors.New("email or username already used")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		userLogger.Warningf("failed to check for existing user: %v", err)
		return err
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		userLogger.Warningf("failed to create user: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		userLogger.Warningf("failed to commit transaction: %v", err)
		return err
	}
	return nil
}

func GetUserByID(id uint32) (*models.User, error) {
	var user models.User
	if err := instance.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := instance.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
