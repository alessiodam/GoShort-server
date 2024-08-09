package v1

import (
	"GoShort/logging"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	routesLogger = logging.NewLogger("internal.router.routes")
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/healthcheck", HealthCheck).Methods("GET")

	r.HandleFunc("/user/register", userRegisterHandler).Methods("POST")
	r.HandleFunc("/user/login", userLoginHandler).Methods("POST")

	r.HandleFunc("/user/me", userMeHandler).Methods("GET")
	r.HandleFunc("/user/logout", userLogoutHandler).Methods("POST")

	r.HandleFunc("/shortlinks", listShortLinksHandler).Methods("GET")
	r.HandleFunc("/shortlinks/{shortlink:[a-zA-Z0-9]+}", getShortLinkHandler).Methods("GET")

	r.HandleFunc("/shortlinks", createShortLinkHandler).Methods("POST")
	r.HandleFunc("/shortlinks/{shortlink:[a-zA-Z0-9]+}", deleteShortLinkHandler).Methods("DELETE")
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		routesLogger.Errorf("Failed to write response: %v", err)
	}
}
