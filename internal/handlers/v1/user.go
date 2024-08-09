package v1

import (
	"GoShort/database"
	"GoShort/internal/handlers/v1/serializers"
	"GoShort/internal/handlers/v1/validators"
	"GoShort/internal/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	requiredFields := []string{"email", "username", "password"}
	for _, field := range requiredFields {
		if _, exists := body[field]; !exists {
			http.Error(w, fmt.Sprintf("Missing required field: %s", field), http.StatusBadRequest)
			return
		}
	}

	email := body["email"].(string)
	if !validators.IsValidEmail(email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	username := body["username"].(string)
	if !validators.IsValidUsername(username) {
		http.Error(w, "Invalid username format", http.StatusBadRequest)
		return
	}

	invalidPassword := validators.IsValidPassword(body["password"].(string))
	if invalidPassword != "" {
		_ = writeJSON(w, http.StatusBadRequest, JSONResponse{Success: false, Message: invalidPassword})
		return
	}

	user := models.User{
		Email:    email,
		Username: username,
	}

	if err := user.HashPassword(body["password"].(string)); err != nil {
		_ = writeJSON(w, http.StatusInternalServerError, JSONResponse{Success: false, Message: err.Error()})
		return
	}

	err := database.InsertUser(&user)
	if err != nil {
		if err.Error() == "email or username already used" {
			_ = writeJSON(w, http.StatusConflict, JSONResponse{Success: false, Message: err.Error()})
		} else {
			_ = writeJSON(w, http.StatusInternalServerError, JSONResponse{Success: false, Message: "Internal Server Error"})
		}
		return
	}

	var newUserSession = models.Session{
		UserID:    user.ID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = database.CreateSession(&newUserSession)
	if err != nil {
		_ = writeJSON(w, http.StatusInternalServerError, JSONResponse{Success: false, Message: "Error creating session"})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "User created successfully",
		"user":    serializers.UserSerializer(&user),
		"session": newUserSession.Token,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		_ = writeJSON(w, http.StatusBadRequest, JSONResponse{Success: false, Message: "Invalid body"})
		return
	}

	requiredFields := []string{"email", "password"}
	for _, field := range requiredFields {
		if _, exists := body[field]; !exists {
			_ = writeJSON(w, http.StatusBadRequest, JSONResponse{Success: false, Message: fmt.Sprintf("Missing required field: %s", field)})
			return
		}
	}

	email := body["email"].(string)
	if !validators.IsValidEmail(email) {
		_ = writeJSON(w, http.StatusBadRequest, JSONResponse{Success: false, Message: "Invalid email format"})
		return
	}

	user, err := database.GetUserByEmail(email)
	if err != nil {
		_ = writeJSON(w, http.StatusNotFound, JSONResponse{Success: false, Message: "User not found"})
		return
	}

	err = models.CheckPassword(user, body["password"].(string))
	if err != nil {
		_ = writeJSON(w, http.StatusBadRequest, JSONResponse{Success: false, Message: "User not found"})
		return
	}

	var newUserSession = models.Session{
		UserID:    user.ID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = database.CreateSession(&newUserSession)
	if err != nil {
		_ = writeJSON(w, http.StatusInternalServerError, JSONResponse{Success: false, Message: "Error creating session"})
		return
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged in successfully!",
		"user":    serializers.UserSerializer(user),
		"session": newUserSession.Token,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func userMeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateWithSessionHeader(r)
	if err != nil {
		_ = writeJSON(w, http.StatusUnauthorized, JSONResponse{Success: false, Message: "Unauthorized"})
		return
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User retrieved successfully",
		"user":    serializers.UserSerializer(&user),
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func userLogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionToken := r.Header.Get("Session")
	if sessionToken == "" {
		_ = writeJSON(w, http.StatusUnauthorized, JSONResponse{Success: false, Message: "Unauthorized"})
		return
	}

	err := database.DeleteSessionByToken(sessionToken)
	if err != nil {
		_ = writeJSON(w, http.StatusInternalServerError, JSONResponse{Success: false, Message: "Error deleting session"})
		return
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged out successfully!",
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
