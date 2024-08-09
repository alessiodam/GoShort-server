package v1

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// writeJSON was taken from this amazing reddit comment
// https://www.reddit.com/r/golang/comments/1egi3i2/comment/lg3hjde/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button
func writeJSON(w http.ResponseWriter, status int, data any) error {
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(js)

	return nil
}
