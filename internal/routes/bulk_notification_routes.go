package routes

import (
	"gorm.io/gorm"
	"github.com/go-chi/chi/v5"

	"github.com/Nitish0007/go_notifier/internal/handlers"
)

func RegisterBulkNotificationRoutes(conn *gorm.DB, r chi.Router, h *handlers.BulkNotificationHandler) {
	r.Post("/bulk_notifications", h.CreateBulkNotificationsHandler) // and return list of errors if any or return successfully created
}