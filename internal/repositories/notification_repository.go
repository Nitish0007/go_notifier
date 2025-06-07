package repositories

import (
	// "context"
	// "github.com/Nitish0007/go_notifier/internal/models"

	"github.com/jackc/pgx/v5"
)

type NotificationRepository struct {
	DB *pgx.Conn
}

func NewNotificationRepository(conn *pgx.Conn) *NotificationRepository {
	return &NotificationRepository{
		DB: conn,
	}
}

// func (r *NotificationRepository) Create(ctx *pgx.Conn, )



