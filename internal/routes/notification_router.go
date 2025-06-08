package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/Nitish0007/go_notifier/internal/handlers"
	"github.com/Nitish0007/go_notifier/internal/middlewares"
)

func RegisterNotificationRoutes(conn *pgx.Conn, r *chi.Mux, h *handlers.NotificationHandler) {
	r.Route("/api/v1/{account_id}", func(r chi.Router) {
		// authenticating request
		r.Use(middlewares.AuthenticateRequest(conn))

		r.Post("/notify", h.SendNotificationHandler)
		r.Post("/bulk_noitify", h.SendBulkNotificationHandler)
	})
}