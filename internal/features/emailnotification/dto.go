package emailnotification

import (
	"time"
)

// Request DTOs
type CreateEmailTransactionalRequest struct {
	Notification struct {
		AccountID        int64  `json:"account_id" validate:"required,gt=0"` // set server-side from auth if omitted
		Title            string `json:"title" validate:"required,max=300"`
		Subject          string `json:"subject" validate:"required,max=500"`
		FromName         string `json:"from_name" validate:"required,max=255"`
		FromEmail        string `json:"from_email" validate:"required,email,max=320"`
		ReplyToEmail     string `json:"reply_to_email" validate:"required,email,max=320"`
		ContentID        int64  `json:"content_id" validate:"required,gt=0"`
		NotificationType string `json:"notification_type" validate:"omitempty,oneof=transactional"`
		Status           string `json:"status" validate:"omitempty,oneof=trans"`
	} `json:"notification" validate:"required"`
}

type CreateEmailCampaignRequest struct {
	Notification struct {
		AccountID        int64   `json:"account_id" validate:"required,gt=0"`
		Title            string  `json:"title" validate:"required,max=300"`
		Subject          string  `json:"subject" validate:"required,max=500"`
		FromName         string  `json:"from_name" validate:"required,max=255"`
		FromEmail        string  `json:"from_email" validate:"required,email,max=320"`
		ReplyToEmail     string  `json:"reply_to_email" validate:"required,email,max=320"`
		ContentID        int64   `json:"content_id" validate:"required,gt=0"`
		ListIDs          []int64  `json:"list_ids" validate:"required"`
		NotificationType string  `json:"notification_type" validate:"omitempty,oneof=campaign"` // defaults to campaign
		Status           string  `json:"status" validate:"required,oneof=draft send_now scheduled"`
		SendAt           *string `json:"send_at" validate:"omitempty"`
	} `json:"notification" validate:"required"`
}

type SendNotificationRequest struct {
	Email    string    `json:"email" validate:"required,email"`
	Subject  string    `json:"subject" validate:"required,max=500"`
	Body     string    `json:"body" validate:"required"`
	HtmlBody string    `json:"html_body" validate:"required"`
	SendAt   time.Time `json:"send_at" validate:"omitempty,datetime"`
}

// Response DTOs
type SendNotificationResponse struct {
	NotificationID string `json:"notification_id"`
}

type EmailCampaignResponse struct {
	ID int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	Subject string `json:"subject"`
	Title string `json:"title"`
	FromName string `json:"from_name"`
	FromEmail string `json:"from_email"`
	ReplyToEmail string `json:"reply_to_email"`
	ContentID int64 `json:"content_id"`
	ListIDs []int64 `json:"list_ids"`
	NotificationType string `json:"notification_type"`
	Status string `json:"status"`
	SendAt *time.Time `json:"send_at"`
	SentAt *time.Time `json:"sent_at"`
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`
}

type CampaignRecipient struct {
	ContactID int64  `json:"contact_id"`
	AccountID int64  `json:"account_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}