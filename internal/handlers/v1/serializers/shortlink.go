package serializers

import "GoShort/internal/models"

func ShortlinkSerializer(shortLink *models.Shortlink) map[string]interface{} {
	return map[string]interface{}{
		"id":        shortLink.ID,
		"long_url":  shortLink.LongURL,
		"short_url": shortLink.ShortURL,
	}
}
