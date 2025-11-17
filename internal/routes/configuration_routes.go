package routes

import (
	"gorm.io/gorm"
	"github.com/go-chi/chi/v5"

	"github.com/Nitish0007/go_notifier/internal/handlers"
)

func RegisterConfigurationRoutes(conn *gorm.DB, r chi.Router, h *handlers.ConfigurationHandler) {
	r.Get("/configurations", h.GetConfigurationsHandler)
	r.Post("/configurations", h.CreateConfigurationHandler)
}