package content

import (
	"gorm.io/gorm"
	"github.com/go-chi/chi/v5"
)

func RegisterContentRoutes(conn *gorm.DB, r chi.Router, h *ContentHandler) {
	r.Post("/contents", h.CreateContentHandler)
}