package list

import (
	"time"
	"github.com/Nitish0007/go_notifier/internal/features/account"
)

type List struct {
	ID            int64 `gorm:"primaryKey" validate:"omitempty,gt=0"`
	AccountID     int64 `gorm:"not null;index" validate:"required,gt=0"`
	Name          string `gorm:"not null" validate:"required,min=1,max=100"`
	ContactsCount int64 `gorm:"not null;default:0" validate:"-"`
	CreatedAt     time.Time `gorm:"autoCreateTime" validate:"-"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" validate:"-"`
	
	// associations
	Account *account.Account `gorm:"foreignKey:AccountID;references:ID" validate:"-"`
}

func NewList(accountId int64, name string) *List {
	return &List{
		AccountID: accountId,
		Name: name,
		ContactsCount: 0,
	}
}