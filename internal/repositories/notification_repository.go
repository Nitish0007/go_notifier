package repositories

import (
	"time"
	"context"

	"github.com/Nitish0007/go_notifier/internal/models"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	DB *gorm.DB
}

func NewNotificationRepository(conn *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		DB: conn,
	}
}

// const allColumns = `
//   id, account_id, channel, recipient, subject, body, html_body, 
//   status, metadata, error_message, job_id, send_at, sent_at, created_at
// `

// func (r *NotificationRepository) Create(ctx context.Context, n *models.Notification) error {

// 	query := `INSERT INTO notifications 
// 	(account_id, channel, recipient, subject, body, html_body, metadata, status, send_at, error_message, job_id) 
// 	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
// 	RETURNING id`

// 	// var jobID any
// 	// if n.JobID == "" {
// 	// 	jobID = nil
// 	// }
// 	err := r.DB.QueryRow(ctx, query, n.AccountID, n.Channel, n.Recipient, n.Subject, n.Body, n.HtmlBody, n.Metadata, n.Status, n.SendAt, n.ErrorMessage, n.JobID).Scan(&n.ID)
// 	if err != nil {
// 		log.Printf("Error creating notification: %v", err)
// 		return err
// 	}
// 	return nil
// }

// func (r *NotificationRepository) Index(ctx context.Context, accID int) ([]*models.Notification, error) {
// 	// query := `SELECT * FROM notifications WHERE account_id = $1 ORDER By created_at DESC`
// 	query := fmt.Sprintf(`SELECT %s FROM notifications WHERE account_id = $1 ORDER BY created_at DESC`, allColumns)
// 	rows, err := r.DB.Query(ctx, query, accID)
// 	if err != nil {
// 		log.Printf("ERROR: %v", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var notifications []*models.Notification

// 	for rows.Next() {
// 		var n models.Notification
// 		err := rows.Scan(
// 			&n.ID,
// 			&n.AccountID,
// 			&n.Channel,
// 			&n.Recipient,
// 			&n.Subject,
// 			&n.Body,
// 			&n.HtmlBody,
// 			&n.Status,
// 			&n.Metadata,
// 			&n.ErrorMessage,
// 			&n.JobID,
// 			&n.SendAt,
// 			&n.SentAt,
// 			&n.CreatedAt,
// 		)
// 		if err != nil {
// 			log.Printf("ERROR: %v", err)
// 			return nil, err
// 		}

// 		notifications = append(notifications, &n)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return notifications, nil
// }

// func (r *NotificationRepository) GetByID(ctx context.Context, id string, accID int) (*models.Notification, error) {
// 	// query := `SELECT * FROM notifications WHERE id = $1 AND account_id = $2`
// 	query := fmt.Sprintf(`SELECT %s FROM notifications WHERE account_id = $1 AND id = $2`, allColumns)

// 	row := r.DB.QueryRow(ctx, query, accID, id)

// 	var n models.Notification
// 	err := row.Scan(
// 		&n.ID,
// 		&n.AccountID,
// 		&n.Channel,
// 		&n.Recipient,
// 		&n.Subject,
// 		&n.Body,
// 		&n.HtmlBody,
// 		&n.Status,
// 		&n.Metadata,
// 		&n.ErrorMessage,
// 		&n.JobID,
// 		&n.SendAt,
// 		&n.SentAt,
// 		&n.CreatedAt,
// 	)
// 	if err != nil {
// 		log.Printf("ERROR: %v", err)
// 		return nil, err
// 	}

// 	return &n, nil
// }


func (r *NotificationRepository) Create(ctx context.Context, n *models.Notification) error {
	return r.DB.WithContext(ctx).Create(n).Error
}

func (r *NotificationRepository) Index(ctx context.Context, accID int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := r.DB.WithContext(ctx).Where("account_id = ?", accID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) GetByID(ctx context.Context, id string, accID int) (*models.Notification, error) {
	var notification models.Notification
	err := r.DB.WithContext(ctx).Where("id = ? AND account_id = ?", id, accID).First(&notification).Error
	return &notification, err
}

func (r *NotificationRepository) UpdateStatus(ctx context.Context, id string, status models.NotificationStatus, errorMsg *string) error {
	updates := map[string]interface{}{"status": status}
	if errorMsg != nil {
			updates["error_message"] = *errorMsg
	}
	if status == models.Sent {
			updates["sent_at"] = time.Now()
	}
	return r.DB.WithContext(ctx).Model(&models.Notification{}).Where("id = ?", id).Updates(updates).Error
}