package repositories

import (
	// "context"
	// "github.com/Nitish0007/go_notifier/internal/models"

	"context"
	"log"

	"github.com/Nitish0007/go_notifier/internal/models"
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

func (r *NotificationRepository) Create(ctx context.Context, n *models.Notification) error {

	query := `INSERT INTO notifications 
	(account_id, channel, recipient, subject, body, html_body, metadata, status, send_at, error_message, job_id) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
	RETURNING id`

	var jobID any
	if n.JobID == "" {
		jobID = nil
	}
	err := r.DB.QueryRow(ctx, query, n.AccountID, n.Channel, n.Recipient, n.Subject, n.Body, n.HtmlBody, n.Metadata, n.Status, n.SendAt, n.ErrorMessage, jobID).Scan(&n.ID)
	if err != nil {
		log.Printf("Error creating notification: %v", err)
		return err
	}
	return nil
}



