package serializers

import (
	"GoShort/database"
	"GoShort/internal/models"
	"GoShort/settings"
)

func ShortlinkSerializer(shortLink *models.Shortlink) map[string]interface{} {
	shortlinkAnalytics, err := database.GetShortlinkAnalyticsByShortlinkID(shortLink.ID)
	if err != nil {
		shortlinkAnalytics = &models.ShortlinkAnalytics{}
	}

	browserAnalytics, err := database.GetShortlinkBrowserAnalyticsByShortlinkID(shortLink.ID)
	if err != nil {
		browserAnalytics = nil
	}

	return map[string]interface{}{
		"id":        shortLink.ID,
		"long_url":  shortLink.LongURL,
		"short_url": settings.Settings.BaseURL + shortLink.ShortURL,
		"analytics": ShortlinkAnalyticsSerializer(shortlinkAnalytics, browserAnalytics),
	}
}

func ShortlinkAnalyticsSerializer(shortLinkAnalytics *models.ShortlinkAnalytics, browserAnalytics []models.ShortlinkBrowserAnalytics) map[string]interface{} {
	return map[string]interface{}{
		"clicks":   shortLinkAnalytics.Clicks,
		"browsers": ShortlinkBrowserAnalyticsSerializer(browserAnalytics),
	}
}

func ShortlinkBrowserAnalyticsSerializer(browserAnalytics []models.ShortlinkBrowserAnalytics) []map[string]interface{} {
	var result []map[string]interface{}
	for _, analytics := range browserAnalytics {
		result = append(result, map[string]interface{}{
			"browser": analytics.Browser,
			"country": analytics.Country,
			"count":   analytics.Count,
		})
	}
	return result
}
