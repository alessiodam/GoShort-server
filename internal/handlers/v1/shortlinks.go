package v1

import (
	"GoShort/database"
	"GoShort/internal/handlers/v1/serializers"
	"GoShort/internal/handlers/v1/validators"
	"GoShort/settings"
	"encoding/json"
	"fmt"
	"net/http"
)

func listShortLinksHandler(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateWithSessionHeader(r)
	if err != nil {
		_ = writeJSON(w, http.StatusUnauthorized, JSONResponse{Success: false, Message: "Unauthorized"})
		return
	}

	shortlinks, err := database.GetShortlinksByUserID(user.ID)
	if err != nil {
		_ = writeJSON(w, http.StatusInternalServerError, JSONResponse{Success: false, Message: "Error retrieving shortlinks"})
		return
	}

	var serializedShortlinks []map[string]interface{}
	for _, shortlink := range shortlinks {
		serializedShortlinks = append(serializedShortlinks, serializers.ShortlinkSerializer(&shortlink))
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"message":    "Shortlinks retrieved successfully",
		"shortlinks": shortlinks,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// TODO
func getShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateWithSessionHeader(r)
	if err != nil {
		_ = writeJSON(w, http.StatusUnauthorized, JSONResponse{Success: false, Message: "Unauthorized"})
		return
	}
}

func createShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateWithSessionHeader(r)
	if err != nil {
		_ = writeJSON(w, http.StatusUnauthorized, JSONResponse{Success: false, Message: "Unauthorized"})
		return
	}

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		_ = writeJSON(w, http.StatusBadRequest, JSONResponse{Success: false, Message: "Invalid body"})
		return
	}

	requiredFields := []string{"url"}
	for _, field := range requiredFields {
		if _, exists := body[field]; !exists {
			_ = writeJSON(w, http.StatusBadRequest, JSONResponse{Success: false, Message: fmt.Sprintf("Missing required field: %s", field)})
			return
		}
	}

	if body["url"] == "" {
		_ = writeJSON(w, http.StatusBadRequest, JSONResponse{Success: false, Message: "Invalid URL"})
		return
	}
	if !validators.IsValidURL(body["url"].(string)) {
		_ = writeJSON(w, http.StatusBadRequest, JSONResponse{Success: false, Message: "Invalid URL (must start with http:// or https://)"})
		return
	}

	shortlink, err := database.CreateShortlink(body["url"].(string), user.ID)
	if err != nil {
		_ = writeJSON(w, http.StatusInternalServerError, JSONResponse{Success: false, Message: "Error creating shortlink"})
		return
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Shortlink created successfully",
		"shortlink": map[string]interface{}{
			"id":        shortlink.ID,
			"long_url":  shortlink.LongURL,
			"short_url": settings.Settings.BaseURL + shortlink.ShortURL,
		},
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// TODO
func deleteShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateWithSessionHeader(r)
	if err != nil {
		_ = writeJSON(w, http.StatusUnauthorized, JSONResponse{Success: false, Message: "Unauthorized"})
		return
	}
}
