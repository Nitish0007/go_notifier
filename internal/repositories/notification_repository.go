package repositories

import (
	// "context"
	// "github.com/Nitish0007/go_notifier/internal/models"

	"context"
	"log"
	"fmt"

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

const allColumns = `
  id, account_id, channel, recipient, subject, body, html_body, 
  status, metadata, error_message, job_id, send_at, sent_at, created_at
`

func (r *NotificationRepository) Create(ctx context.Context, n *models.Notification) error {

	query := `INSERT INTO notifications 
	(account_id, channel, recipient, subject, body, html_body, metadata, status, send_at, error_message, job_id) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
	RETURNING id`

	// var jobID any
	// if n.JobID == "" {
	// 	jobID = nil
	// }
	err := r.DB.QueryRow(ctx, query, n.AccountID, n.Channel, n.Recipient, n.Subject, n.Body, n.HtmlBody, n.Metadata, n.Status, n.SendAt, n.ErrorMessage, n.JobID).Scan(&n.ID)
	if err != nil {
		log.Printf("Error creating notification: %v", err)
		return err
	}
	return nil
}

func (r *NotificationRepository) Index(ctx context.Context, accID int) ([]*models.Notification, error) {
	// query := `SELECT * FROM notifications WHERE account_id = $1 ORDER By created_at DESC`
	query := fmt.Sprintf(`SELECT %s FROM notifications WHERE account_id = $1 ORDER BY created_at DESC`, allColumns)
	rows, err := r.DB.Query(ctx, query, accID)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil, err
	}
	defer rows.Close()

	var notifications []*models.Notification

	for rows.Next() {
		var n models.Notification
		err := rows.Scan(
			&n.ID,
			&n.AccountID,
			&n.Channel,
			&n.Recipient,
			&n.Subject,
			&n.Body,
			&n.HtmlBody,
			&n.Status,
			&n.Metadata,
			&n.ErrorMessage,
			&n.JobID,
			&n.SendAt,
			&n.SentAt,
			&n.CreatedAt,
		)
		if err != nil {
			log.Printf("ERROR: %v", err)
			return nil, err
		}

		notifications = append(notifications, &n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *NotificationRepository) GetByID(ctx context.Context, id string, accID int) (*models.Notification, error) {
	// query := `SELECT * FROM notifications WHERE id = $1 AND account_id = $2`
	query := fmt.Sprintf(`SELECT %s FROM notifications WHERE account_id = $1 AND id = $2`, allColumns)

	row := r.DB.QueryRow(ctx, query, accID, id)

	var n models.Notification
	err := row.Scan(
		&n.ID,
		&n.AccountID,
		&n.Channel,
		&n.Recipient,
		&n.Subject,
		&n.Body,
		&n.HtmlBody,
		&n.Status,
		&n.Metadata,
		&n.ErrorMessage,
		&n.JobID,
		&n.SendAt,
		&n.SentAt,
		&n.CreatedAt,
	)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return nil, err
	}

	return &n, nil
}