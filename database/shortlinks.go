package database

import (
	"GoShort/internal/models"
	"math/rand"
	"time"
)

func GetShortlinksByUserID(userID uint32) ([]models.Shortlink, error) {
	var shortlinks []models.Shortlink
	if err := instance.Where("user_id = ?", userID).Find(&shortlinks).Error; err != nil {
		return nil, err
	}
	return shortlinks, nil
}

func GetShortlinkByShortURL(shortURL string) (*models.Shortlink, error) {
	var shortlink models.Shortlink
	if err := instance.Where("short_url = ?", shortURL).First(&shortlink).Error; err != nil {
		return nil, err
	}
	return &shortlink, nil
}

func generateShortURL() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func CreateShortlink(originalURL string, userID uint32) (models.Shortlink, error) {
	shortURL := generateShortURL()

	shortlink := models.Shortlink{
		ShortURL: shortURL,
		LongURL:  originalURL,
		UserID:   userID,
	}

	if err := instance.Create(&shortlink).Error; err != nil {
		return models.Shortlink{}, err
	}

	var newShortlink models.Shortlink
	if err := instance.First(&newShortlink, shortlink.ID).Error; err != nil {
		return models.Shortlink{}, err
	}

	return newShortlink, nil
}
