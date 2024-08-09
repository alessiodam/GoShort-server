package router

import (
	"GoShort/internal/handlers/v1"
	"GoShort/internal/middlewares"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(middlewares.CORSMiddleware)
	r.Use(middlewares.JSONMiddleware)

	apiV1Router := r.PathPrefix("/api/v1").Subrouter()
	v1.RegisterRoutes(apiV1Router)

	r.HandleFunc("/", indexHandler)
	shortlinkRouter := r.PathPrefix("/").Subrouter()
	shortlinkRouter.HandleFunc("/{shortlink:[a-zA-Z0-9]+}", handleShortlink)

	return r
}
