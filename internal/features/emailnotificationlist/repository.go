package emailnotificationlist

import (
	// "context"
	"gorm.io/gorm"
)

type EmailNotificationListRepository struct {
	db *gorm.DB
}

func NewEmailNotificationListRepository(db *gorm.DB) *EmailNotificationListRepository {
	return &EmailNotificationListRepository{db: db}
}