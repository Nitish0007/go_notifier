package models

import (
	"errors"
	"time"
)

type NotificationStatus int
type NotificationChannel int

const (
	Email NotificationChannel = iota
	Sms
	InApp
)

const (
	Pending	NotificationStatus = iota
	Enqueued
	Sent
	failed
)

type Notification struct {
	ID           string               `json:"id" gorm:"type:uuid;primaryKey;"`
	AccountID    int                  `json:"account_id" gorm:"not null;index"`
	Channel      NotificationChannel  `json:"channel" gorm:"not null;check:channel IN (0,1,2)"`
	Recipient    string               `json:"recipient" gorm:"not null;size:255"`
	Subject      string               `json:"subject" gorm:"size:500"`
	Body         string               `json:"body" gorm:"type:text"`
	HtmlBody     string               `json:"html_body" gorm:"type:text"`
	Status       NotificationStatus   `json:"status" gorm:"not null;default:0;check:status IN (0,1,2,3)"`
	Metadata     map[string]any       `json:"metadata" gorm:"type:jsonb;default:'{}'"`
	ErrorMessage *string              `json:"error_message" gorm:"type:text"`
	JobID        *string              `json:"job_id" gorm:"type:uuid"`
	SendAt       *time.Time           `json:"send_at"`
	SentAt       *time.Time           `json:"sent_at"`
	CreatedAt    time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time            `json:"updated_at" gorm:"autoUpdateTime"`
	BatchID      *string              `json:"batch_id" gorm:"type:uuid;index"`
}

func StringToNotificationStatus(status string) (NotificationStatus, error) {
	switch status {
	case "pending":
		return Pending, nil
	case "enqueued":
		return Enqueued, nil
	case "sent":
		return Sent, nil
	case "failed":
		return failed, nil
	default:
		return -1, errors.New("unknown status type")
	}
}

func StatusToString(status NotificationStatus) (string, error ){
	switch status {
	case Pending:
		return "pending", nil
	case Enqueued:
		return "enqueued", nil
	case Sent:
		return "sent", nil
	case failed:
		return "failed", nil
	default:
		return "", errors.New("unknown status type")
	}
}

func StringToNotificationChannel(channel string) (NotificationChannel, error) {
	switch channel {
	case "email":
		return Email, nil
	case "sms":
		return Sms, nil
	case "in_app":
		return InApp, nil
	default:
		return -1, errors.New("unknown channel type")
	}
}

func ChannelToString(channel NotificationChannel) (string, error) {
	switch channel {
	case Email:
		return "email", nil
	case Sms:
		return "sms", nil
	case InApp:
		return "in_app", nil
	default:
		return "", errors.New("unknown channel type")
	}
}