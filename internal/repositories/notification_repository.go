package repositories

import (
	"context"

	"github.com/Nitish0007/go_notifier/internal/models"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	DB *gorm.DB
}

func NewNotificationRepository(conn *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		DB: conn,
	}
}

func (r *NotificationRepository) Create(ctx context.Context, n *models.Notification) error {
	return r.DB.WithContext(ctx).Create(n).Error
}

func (r *NotificationRepository) Index(ctx context.Context, accID int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := r.DB.WithContext(ctx).Where("account_id = ?", accID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) GetByID(ctx context.Context, id string, accID int) (*models.Notification, error) {
	var notification models.Notification
	err := r.DB.WithContext(ctx).Where("id = ? AND account_id = ?", id, accID).First(&notification).Error
	return &notification, err
}