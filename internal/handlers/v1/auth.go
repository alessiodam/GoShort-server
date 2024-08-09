package v1

import (
	"GoShort/database"
	"GoShort/internal/models"
	"errors"
	"net/http"
)

func authenticateWithSessionHeader(r *http.Request) (models.User, error) {
	sessionHeader := r.Header.Get("Session")
	if sessionHeader == "" {
		return models.User{}, errors.New("no Session header :(")
	}
	user, err := database.GetUserBySessionToken(sessionHeader)
	if err != nil {
		return models.User{}, errors.New("user no exist :(")
	}
	return *user, nil
}
