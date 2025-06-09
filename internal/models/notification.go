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
	ID						string								`json:"id"`
	AccountID 		int										`json:"account_id"`
	Channel				NotificationChannel		`json:"channel"`
	Recipient			string								`json:"recipient"`
	Subject				string								`json:"subject"`
	Body					string								`json:"body"`
	HtmlBody			string								`json:"html_body"`
	Status				NotificationStatus		`json:"status"`
	Metadata			map[string]any				`json:"metadata"`
	ErrorMessage	string								`json:"error_message"`
	JobID					string								`json:"job_id"`
	SendAt				time.Time							`json:"send_at"`
	SentAt				time.Time							`json:"sent_at"`
	CreatedAt			time.Time							`json:"created_at"`
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