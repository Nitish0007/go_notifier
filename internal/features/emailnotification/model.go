package emailnotification

import (
	"time"
	"errors"
	"strings"
	"encoding/json"
)

type EmailNotificationStatus int
type EmailNotificationType int

const (
	Transactional EmailNotificationType = iota
	Campaign
)

const (
	Trans EmailNotificationStatus = iota
	Draft
	Scheduled
	Enqueued
	Sent
	Failed
)

type EmailNotification struct {
	ID               int64                    `json:"id" gorm:"primaryKey" validate:"omitempty,gt=0"`
	AccountID        int                      `json:"account_id" gorm:"not null;index" validate:"required,gt=0"`
	Subject          string                   `json:"subject" gorm:"size:500" validate:"omitempty,max=500"`
	Title            string                   `json:"title" gorm:"size:300" validate:"omitempty,max=300"`
	NotificationType EmailNotificationType    `json:"notification_type" gorm:"not null;default:0;check:notification_type IN (0,1)" validate:"-"` // 0 = transactional, 1 = campaign
	ContentID        int                      `json:"content_id" gorm:"not null;index" validate:"omitempty,gt=0"`
	Status           EmailNotificationStatus  `json:"status" gorm:"not null;default:0;check:status IN (0,1,2,3)" validate:"-"` // Custom validation needed for enum
	SendAt           *time.Time               `json:"send_at" validate:"-"`
	SentAt           *time.Time               `json:"sent_at" validate:"-"`
	CreatedAt        time.Time                `json:"created_at" gorm:"autoCreateTime" validate:"-"`
}

func NewEmailNotification(accountId int, subject string, title string, notificationType EmailNotificationType, contentId int, status EmailNotificationStatus, sendAt *time.Time) *EmailNotification {
	return &EmailNotification{
		AccountID: accountId,
		Subject: subject,
		Title: title,
		NotificationType: notificationType,
		ContentID: contentId,
		Status: status,
		SendAt: sendAt,
		SentAt: nil,
	}
}

func StringToEmailNotificationType(notificationType string) (EmailNotificationType, error) {
	typeStr := strings.ToLower(notificationType)
	switch typeStr {
	case "transactional":
		return Transactional, nil
	case "campaign":
		return Campaign, nil
	default:
		return 0, errors.New("invalid notification type")
	}
}

func StringToEmailNotificationStatus(status string) (EmailNotificationStatus, error) {
	switch status {
	case "trans":
		return Trans, nil
	case "draft":
		return Draft, nil
	case "scheduled":
		return Scheduled, nil
	case "enqueued":
		return Enqueued, nil
	case "sent":
		return Sent, nil
	case "failed":
		return Failed, nil
	default:
		return -1, errors.New("invalid status type")
	}
}

func StatusToString(status EmailNotificationStatus) (string, error) {
	switch status {
	case Trans:
		return "pending", nil
	case Draft:
		return "draft", nil
	case Scheduled:
		return "scheduled", nil
	case Enqueued:
		return "enqueued", nil
	case Sent:
		return "sent", nil	
	case Failed:
		return "failed", nil
	default:
		return "", errors.New("unknown status type")
	}
}

// to convert notification struct to map[string]any
func (n *EmailNotification) ToMap() (map[string]any, error) {
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

func (n *EmailNotification) ResponseMap() (map[string]any, error) {
	statusString, err := StatusToString(n.Status)
	if err != nil {
		return nil, err	
	}
	return map[string]any{
		"id": n.ID,
		"account_id": n.AccountID,
		"subject": n.Subject,
		"title": n.Title,
		"notification_type": n.NotificationType,
		"content_id": n.ContentID,
		"status": statusString,
		"send_at": n.SendAt,
		"sent_at": n.SentAt,
		"created_at": n.CreatedAt,
	}, nil
}


// Implement NotificationView (getter names avoid conflict with struct fields).
func (n *EmailNotification) GetSubject() string     					{ return n.Subject }
func (n *EmailNotification) GetTitle() string     					{ return n.Title }
func (n *EmailNotification) GetNotificationType() EmailNotificationType 	{ return n.NotificationType }
func (n *EmailNotification) GetContentID() int 						{ return n.ContentID }
// func (n *EmailNotification) GetStatus() EmailNotificationStatus 	{ return n.Status }
// func (n *EmailNotification) GetAccountID() int 							{ return n.AccountID }
// func (n *EmailNotification) GetSendAt() *time.Time 					{ return n.SendAt }
// func (n *EmailNotification) GetSentAt() *time.Time 					{ return n.SentAt }
// func (n *EmailNotification) GetCreatedAt() time.Time 				{ return n.CreatedAt }
// func (n *EmailNotification) GetID() string 									{ return n.ID }