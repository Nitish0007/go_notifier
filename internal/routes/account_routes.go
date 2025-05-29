package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/Nitish0007/go_notifier/internal/handlers"
	"github.com/Nitish0007/go_notifier/internal/middlewares"
)

func RegisterAccountRoutes(r *chi.Mux, h *handlers.AccountHandler){
	// public routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/signup", h.CreateAccountHandler)
		r.Get("/login", h.LoginHandler)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middlewares.AuthenticateRequest)

		// r.Get("/accounts", h.GetAccountsHandler)
	})

}