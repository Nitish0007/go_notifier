package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/Nitish0007/go_notifier/internal/handlers"
)

func RegisterPublicAccountRoutes(conn *pgx.Conn, r chi.Router, h *handlers.AccountHandler) {
	// public routes
	r.Post("/signup", h.CreateAccountHandler)
	r.Get("/login", h.LoginHandler)
}

func RegisterAccountRoutes(conn *pgx.Conn, r chi.Router, h *handlers.AccountHandler){
	// protected routes
 // any route defined here will be authenticated
}