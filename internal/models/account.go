package models

import (
	"time"
)

type Account struct {
	ID                int       `json:"id" gorm:"primaryKey" validate:"-"`
	Email             string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	EncryptedPassword string    `json:"-" gorm:"column:encrypted_password;not null" validate:"required,min=8"` // do not expose in JSON
	FirstName         string    `json:"first_name" gorm:"column:first_name;not null" validate:"required,min=1,max=100"`
	LastName          string    `json:"last_name" gorm:"column:last_name;not null" validate:"required,min=1,max=100"`
	IsActive          bool      `json:"is_active" gorm:"column:is_active;default:true" validate:"-"`
	CreatedAt         time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime" validate:"-"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime" validate:"-"`
}
