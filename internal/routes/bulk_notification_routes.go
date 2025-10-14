package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Nitish0007/go_notifier/internal/handlers"
)

func RegisterBulkNotificationRoutes(conn *pgxpool.Pool, r chi.Router, h *handlers.BulkNotificationHandler) {
	r.Post("/bulk_notifications/", h.CreateBulkNotificationsHandler) // and return list of errors if any or return successfully created
}