package models

import (
	"time"
)

type Account struct {
	ID                int       `json:"id" gorm:"primaryKey"`
	Email             string    `json:"email" gorm:"uniqueIndex;not null"`
	EncryptedPassword string    `json:"-" gorm:"column:encrypted_password;not null"` // do not expose in JSON
	FirstName         string    `json:"first_name" gorm:"column:first_name;not null"`
	LastName          string    `json:"last_name" gorm:"column:last_name;not null"`
	IsActive          bool      `json:"is_active" gorm:"column:is_active;default:true"`
	CreatedAt         time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}