package notification

import (
	"time"
)

// Request DTOs
type SendNotificationRequest struct {
	Channel string `json:"channel" validate:"required,oneof=email sms in_app"`
	Recipient string `json:"recipient" validate:"required,email"`
	Subject string `json:"subject" validate:"required,max=500"`
	Body string `json:"body" validate:"required"`
	HtmlBody string `json:"html_body" validate:"required"`
	SendAt time.Time `json:"send_at" validate:"omitempty,datetime"`
}

// Response DTOs
type SendNotificationResponse struct {
	NotificationID string `json:"notification_id"`
}