package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationStatus int
type NotificationChannel int

const (
	Email NotificationChannel = iota
	Sms
	InApp
)

const (
	Pending NotificationStatus = iota
	Enqueued
	Sent
	failed
)

type Notification struct {
	ID           string              `json:"id" gorm:"type:uuid;primaryKey;" validate:"omitempty,uuid"`
	AccountID    int                 `json:"account_id" gorm:"not null;index" validate:"required,gt=0"`
	Channel      NotificationChannel `json:"channel" gorm:"not null;check:channel IN (0,1,2)" validate:"required"` // Custom validation needed for enum
	Recipient    string              `json:"recipient" gorm:"not null;size:255" validate:"required,min=1,max=255"`
	Subject      string              `json:"subject" gorm:"size:500" validate:"omitempty,max=500"`
	Body         string              `json:"body" gorm:"type:text" validate:"omitempty"`
	HtmlBody     string              `json:"html_body" gorm:"type:text" validate:"omitempty"`
	Status       NotificationStatus  `json:"status" gorm:"not null;default:0;check:status IN (0,1,2,3)" validate:"-"` // Custom validation needed for enum
	Metadata     map[string]any      `json:"metadata" gorm:"type:jsonb; default:'{}'; serializer:json" validate:"-"`
	ErrorMessage *string             `json:"error_message" gorm:"type:text" validate:"-"`
	JobID        *string             `json:"job_id" gorm:"type:uuid" validate:"omitempty,uuid"`
	SendAt       *time.Time          `json:"send_at" validate:"-"`
	SentAt       *time.Time          `json:"sent_at" validate:"-"`
	CreatedAt    time.Time           `json:"created_at" gorm:"autoCreateTime" validate:"-"`
	BatchID      *string             `json:"batch_id" gorm:"type:uuid;index" validate:"omitempty,uuid"`
}

// Before Create hook to generate UUID
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return nil
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

func StatusToString(status NotificationStatus) (string, error) {
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

// to convert notification struct to map[string]any
func (n *Notification) ToMap() (map[string]any, error) {
	// marshal struct to JSON
	jsonBytes, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON to ma[string]any
	var result map[string]any
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}
