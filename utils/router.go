package utils

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// initializing and returning router with basic middlewares and configurations
func InitRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(middleware.Timeout(60 * time.Second)) // Set a timeout of 60 seconds for requests
	router.Use(middleware.URLFormat)                 // Parse extension from url and put it on request context, like .json, .xml
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300, // Maximum value not ignored by any of major browsers
	}))

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Go Notifier API!"))

		panic("Simulated panic for testing purpose")
	})

	return router
}
