package contact

import (
	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/emailcontact"
	"time"
)

type Contact struct {
	ID        int       `json:"id" gorm:"primaryKey" validate:"omitempty,gt=0"`
	UUID      string    `json:"uuid" gorm:"uniqueIndex" validate:"omitempty,uuid"`
	AccountID int       `json:"account_id" gorm:"not null;index" validate:"required,gt=0"`
	FirstName string    `json:"first_name" gorm:"not null" validate:"required,min=1,max=100"`
	LastName  string    `json:"last_name" gorm:"not null" validate:"required,min=1,max=100"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" validate:"-"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" validate:"-"`

	// associations
	Account      *account.Account           `json:"-" gorm:"foreignKey:AccountID;references:ID" validate:"-"`
	EmailContact *emailcontact.EmailContact `json:"-" gorm:"foreignKey:ContactID;references:ID" validate:"-"`
}
