package models

import (
	"time"
	"errors"
	"gorm.io/gorm"
)

type NotificationBatch struct {
	ID              string     			`json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AccountID       int        			`json:"account_id" gorm:"not null;index"`
	Payload         map[string]any 	`json:"payload" gorm:"type:jsonb;not null"` // payload of the notifications
	Count           int        			`json:"count" gorm:"not null"`
	SuccessfulCount int        			`json:"successful_count" gorm:"default:0"` // count of successful notifications
	FailedCount     int        			`json:"failed_count" gorm:"default:0"` // count of failed notifications
	Channel         int        			`json:"channel" gorm:"not null;check:channel IN (0,1,2)"`
	Status          int        			`json:"status" gorm:"not null;default:0;check:status IN (0,1,2,3)"` // [0 - pending, 1 - enqueued, 2 - sent, 3 - failed]
	CreatedAt       time.Time  			`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  			`json:"updated_at" gorm:"autoUpdateTime"`
	CompletedAt     *time.Time 			`json:"completed_at" gorm:"autoUpdateTime;default:null"` // timestamp when batch is completed
}

func (nb *NotificationBatch) BeforeSave(tx *gorm.DB) error {
	if len(nb.Payload) == 0 {
		return errors.New("payload is required")
	}
	return nil
}