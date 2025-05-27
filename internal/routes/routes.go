package routes

import (
	"github.com/Nitish0007/go_notifier/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func InitializeRoutes(r *chi.Mux){
	// Accounts routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/accounts", handlers.CreateAccountHandler)
	})
}