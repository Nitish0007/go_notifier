package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Nitish0007/go_notifier/internal/handlers"
)

func RegisterNotificationRoutes(conn *pgxpool.Pool, r chi.Router, h *handlers.NotificationHandler) {
	r.Post("/notify", h.SendNotificationHandler) // create and send notification and return the notification id
	// r.Post("/bulk_notify", h.SendBulkNotificationHandler) 
	r.Post("/notification", h.CreateNotificationHandler) // create notification and return its object
	r.Post("/trigger", h.SendNotificationByIDHandler) // send notification by id
	r.Get("/notifications", h.GetNotificationsHandler) // get notifications in context of account
}