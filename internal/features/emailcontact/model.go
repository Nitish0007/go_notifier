package emailcontact

import (
	"time"
)

type EmailContact struct {
	ID        int       `json:"id" gorm:"primaryKey" validate:"omitempty,gt=0"`
	AccountID int       `json:"account_id" gorm:"not null;index" validate:"required,gt=0"`
	ContactID int       `json:"contact_id" gorm:"not null;index" validate:"required,gt=0"`
	Email     string    `json:"email" gorm:"not null;index" validate:"required,email"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" validate:"-"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" validate:"-"`
}
