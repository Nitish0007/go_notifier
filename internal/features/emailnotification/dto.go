package emailnotification

import (
	"time"
)

// Request DTOs
type CreateEmailTransactionalRequest struct {
	Notification struct {
		AccountID  		   int 	                   `json:"account_id" validate:"required,gt=0"`
		Title            string                  `json:"title" validate:"required,max=300"`
		Subject          string                  `json:"subject" validate:"required,max=500"`
		ContentID        int                     `json:"content_id" validate:"omitempty,gt=0"`
		NotificationType string                  `json:"notification_type" validate:"omitempty,oneof=transactional"`
		Status           string                  `json:"status" validate:"omitempty,oneof=trans"`
	} `json:"notification" validate:"required"`
}

type CreateEmailCampaignRequest struct {
	Notification struct {
		AccountID  		   int 	                   `json:"account_id" validate:"required,gt=0"`
		Title            string                  `json:"title" validate:"required,max=300"`
		Subject          string                  `json:"subject" validate:"required,max=500"`
		ContentID        int                     `json:"content_id" validate:"omitempty,gt=0"`
		NotificationType string                  `json:"notification_type" validate:"omitempty,oneof=campaign"`
		Status           string                  `json:"status" validate:"required,oneof=trans draft send_now scheduled"`
		SendAt           *string                 `json:"send_at" validate:"omitempty"`
	} `json:"notification" validate:"required"`
}

type SendNotificationRequest struct {
	Email string `json:"email" validate:"required,email"`
	Subject string `json:"subject" validate:"required,max=500"`
	Body string `json:"body" validate:"required"`
	HtmlBody string `json:"html_body" validate:"required"`
	SendAt time.Time `json:"send_at" validate:"omitempty,datetime"`
}

// Response DTOs
type SendNotificationResponse struct {
	NotificationID string `json:"notification_id"`
}