package configuration

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterConfigurationRoutes(conn *gorm.DB, r chi.Router, h *ConfigurationHandler) {
	r.Get("/configurations", h.GetConfigurationsHandler)
	r.Post("/configurations", h.CreateConfigurationHandler)
	r.Delete("/configurations/{id}", h.DeleteConfigurationHandler)
	r.Put("/configurations/{id}", h.UpdateConfigurationHandler)
	r.Patch("/configurations/{id}", h.UpdateConfigurationHandler)
}
