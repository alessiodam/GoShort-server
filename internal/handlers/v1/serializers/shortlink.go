package serializers

import (
	"GoShort/database"
	"GoShort/internal/models"
	"GoShort/settings"
)

func ShortlinkSerializer(shortLink *models.Shortlink) map[string]interface{} {
	shortlinkAnalytics, err := database.GetShortlinkAnalyticsByShortlinkID(shortLink.ID)
	if err != nil {
		return map[string]interface{}{
			"id":        shortLink.ID,
			"long_url":  shortLink.LongURL,
			"short_url": settings.Settings.BaseURL + shortLink.ShortURL,
			"analytics": nil,
		}
	}
	return map[string]interface{}{
		"id":        shortLink.ID,
		"long_url":  shortLink.LongURL,
		"short_url": settings.Settings.BaseURL + shortLink.ShortURL,
		"analytics": ShortlinkAnalyticsSerializer(shortlinkAnalytics),
	}
}

func ShortlinkAnalyticsSerializer(shortLinkAnalytics *models.ShortlinkAnalytics) map[string]interface{} {
	return map[string]interface{}{
		"clicks": shortLinkAnalytics.Clicks,
	}
}
