package models

import (
	"time"
)

type ApiKey struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"not null;unique"`
	AccountID int       `json:"account_id" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
