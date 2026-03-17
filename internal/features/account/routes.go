package account

import (
	"gorm.io/gorm"
	"github.com/go-chi/chi/v5"
)

func RegisterPublicAccountRoutes(conn *gorm.DB, r chi.Router, h *AccountHandler) {
	// public routes
	r.Post("/signup", h.SignupHandler)
	r.Post("/login", h.LoginHandler)
}

func RegisterAccountRoutes(conn *gorm.DB, r chi.Router, h *AccountHandler){
	// protected routes
 // any route defined here will be authenticated
}