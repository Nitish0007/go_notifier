package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/Nitish0007/go_notifier/internal/handlers"
)

func RegisterAccountRoutes(r *chi.Mux, h *handlers.AccountHandler){
	// Accounts routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/accounts", h.CreateAccountHandler)
	})
}