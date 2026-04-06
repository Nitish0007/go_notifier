package emailnotificationlist

import (
	"time"
)

type EmailNotificationList struct {
	ID               int64     `json:"id" gorm:"primaryKey" validate:"omitempty,gt=0"`
	AccountID        int64     `json:"account_id" gorm:"not null;index" validate:"required,gt=0"`
	ListID           int64     `json:"list_id" gorm:"not null;index" validate:"required,gt=0"`
	NotificationID   int64     `json:"notification_id" gorm:"not null;index" validate:"required,gt=0"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime" validate:"-"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime" validate:"-"`
}