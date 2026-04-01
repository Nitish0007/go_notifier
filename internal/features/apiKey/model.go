package apiKey

import (
	"time"
)

// ApiKey is an API credential row for an account (one active key per account in current schema).
type ApiKey struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"not null;unique"`
	AccountID int       `json:"account_id" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
