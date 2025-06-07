package models

import (
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
	Content				string								`json:"Content"`
	Status				NotificationStatus		`json:"status"`
	Metadata			map[string]any				`json:"metadata"`
	ErrorMessage	string								`json:"error_message"`
	JobID					string								`json:"job_id"`
	SendAt				time.Time							`json:"send_at"`
	SentAt				time.Time							`json:"sent_at"`
	CreatedAt			time.Time							`json:"created_at"`
}