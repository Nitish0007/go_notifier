package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/Nitish0007/go_notifier/internal/handlers"
)

func RegisterNotificationRoutes(conn *pgx.Conn, r chi.Router, h *handlers.NotificationHandler) {
	r.Post("/notify", h.SendNotificationHandler)
	r.Post("/bulk_notify", h.SendBulkNotificationHandler)
	r.Post("/trigger", h.SendNotificationByIDHandler)
	r.Get("/notifications", h.GetNotificationsHandler)
}