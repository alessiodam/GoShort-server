package database

import (
	"GoShort/internal/models"
	"errors"
	"gorm.io/gorm"
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

	shortlinkAnalytics := models.ShortlinkAnalytics{
		ShortlinkID: newShortlink.ID,
		Clicks:      0,
	}
	if err := instance.Create(&shortlinkAnalytics).Error; err != nil {
		return models.Shortlink{}, err
	}

	return newShortlink, nil
}

func GetShortlinkAnalyticsByShortlinkID(shortlinkID uint32) (*models.ShortlinkAnalytics, error) {
	var shortlinkAnalytics models.ShortlinkAnalytics
	if err := instance.Where("shortlink_id = ?", shortlinkID).First(&shortlinkAnalytics).Error; err != nil {
		return nil, err
	}
	return &shortlinkAnalytics, nil
}

func RecordClickWithDetails(shortlinkID uint32, browser, country string) error {
	var analytics models.ShortlinkAnalytics
	instance.Where("shortlink_id = ?", shortlinkID).First(&analytics)
	analytics.Clicks++
	if err := instance.Save(&analytics).Error; err != nil {
		return err
	}

	var browserAnalytics models.ShortlinkBrowserAnalytics
	result := instance.Where("shortlink_id = ? AND browser = ? AND country = ?", shortlinkID, browser, country).First(&browserAnalytics)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		browserAnalytics = models.ShortlinkBrowserAnalytics{
			ShortlinkID: shortlinkID,
			Browser:     browser,
			Country:     country,
			Count:       1,
		}
		if err := instance.Create(&browserAnalytics).Error; err != nil {
			return err
		}
	} else {
		browserAnalytics.Count++
		if err := instance.Save(&browserAnalytics).Error; err != nil {
			return err
		}
	}

	return nil
}

func GetShortlinkBrowserAnalyticsByShortlinkID(shortlinkID uint32) ([]models.ShortlinkBrowserAnalytics, error) {
	var analytics []models.ShortlinkBrowserAnalytics
	err := instance.Where("shortlink_id = ?", shortlinkID).Find(&analytics).Error
	if err != nil {
		return nil, err
	}
	return analytics, nil
}
