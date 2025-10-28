package models

import (
	"time"
)

type NotificationBatchError struct {
	ID           string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BatchID      string         `json:"batch_id" gorm:"not null;index"`
	BatchIndex   int            `json:"batch_index" gorm:"not null"`
	AccountID    int            `json:"account_id" gorm:"not null;index"`
	ErrorMessage string         `json:"error_message" gorm:"not null;type:text"`
	ErrorType    string         `json:"error_type" gorm:"size:50;default:'validation'"`
	Payload      map[string]any `json:"payload" gorm:"type:jsonb"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}