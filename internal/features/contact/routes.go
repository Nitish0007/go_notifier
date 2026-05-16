package contact

import (
	"gorm.io/gorm"
	"github.com/go-chi/chi/v5"
)

func RegisterContactRoutes(conn *gorm.DB, r chi.Router, h *ContactHandler) {
	r.Get("/contacts", h.GetContactsHandler)
	r.Post("/contacts", h.CreateContactHandler)
	r.Get("/contact/{id}", h.GetContactByIdHandler)
	// r.Put("/contacts/{id}", h.UpdateContactHandler)
	// r.Delete("/contacts/{id}", h.DeleteContactHandler) // soft deletion only
}