package routes

import (
	"gorm.io/gorm"
	"github.com/go-chi/chi/v5"

	"github.com/Nitish0007/go_notifier/internal/handlers"
)

func RegisterPublicAccountRoutes(conn *gorm.DB, r chi.Router, h *handlers.AccountHandler) {
	// public routes
	r.Post("/signup", h.CreateAccountHandler)
	r.Post("/login", h.LoginHandler)
}

func RegisterAccountRoutes(conn *gorm.DB, r chi.Router, h *handlers.AccountHandler){
	// protected routes
 // any route defined here will be authenticated
}