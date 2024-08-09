package serializers

import "GoShort/internal/models"

func UserSerializer(user *models.User) map[string]interface{} {
	return map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	}
}
