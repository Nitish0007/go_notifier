package models

import (
	"time"
)

type NotificationBatch struct {
	ID              string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AccountID       int        `json:"account_id" gorm:"not null;index"`
	BatchID         string     `json:"batch_id" gorm:"type:uuid;not null;uniqueIndex"`
	TotalRequested  int        `json:"total_requested" gorm:"not null"`
	TotalProcessed  int        `json:"total_processed" gorm:"default:0"`
	SuccessfulCount int        `json:"successful_count" gorm:"default:0"`
	FailedCount     int        `json:"failed_count" gorm:"default:0"`
	Status          string     `json:"status" gorm:"not null;default:'processing';check:status IN ('processing','completed','failed')"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	CompletedAt     *time.Time `json:"completed_at" gorm:"autoUpdateTime"`
}
