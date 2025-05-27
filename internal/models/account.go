package models

import (
	"time"
)

type Account struct {
	ID                int       `json:"id" db:"id"`
	Email             string    `json:"email" db:"email"`
	EncryptedPassword string    `json:"-" db:"encrypted_password"` // do not expose in JSON
	FirstName         string    `json:"first_name" db:"first_name"`
	LastName          string    `json:"last_name" db:"last_name"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}