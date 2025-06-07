package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/Nitish0007/go_notifier/internal/handlers"
)

func RegisterNotificationRoutes(conn *pgx.Conn, r *chi.Mux, h *handlers.NotificationHandler) {
	
}