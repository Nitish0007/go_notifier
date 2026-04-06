package listsubscription

import (
	"time"
)

type ListSubscription struct {
	ID        int64     `gorm:"primaryKey" validate:"omitempty,gt=0"`
	AccountID int64     `gorm:"not null;index" validate:"required,gt=0"`
	ListID    int64     `gorm:"not null;index" validate:"required,gt=0"`
	ContactID int64     `gorm:"not null;index" validate:"required,gt=0"`
	Active    bool      `gorm:"not null;default:false" validate:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" validate:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" validate:"-"`
}

func NewListSubscription(accountId int64, listId int64, contactId int64, active bool) *ListSubscription {
	return &ListSubscription{
		AccountID: accountId,
		ListID:    listId,
		ContactID: contactId,
		Active:    active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}