package repositories

import (
	"gorm.io/gorm"
)

type NotificationBatchErrorRepo struct {
	DB *gorm.DB
}

func NewNotificationBatchErrorRepo(conn *gorm.DB) *NotificationBatchErrorRepo {
	return &NotificationBatchErrorRepo{
		DB: conn,
	}
}