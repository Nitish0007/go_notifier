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

func (r *NotificationRepository) GetNotificationsByChannel(ctx context.Context, ch string, accID int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	channel, err := models.StringToNotificationChannel(ch)
	if err != nil {
		return nil, err
	}

	n := &models.Notification{
		AccountID: accID,
		Channel: channel,
	}

	// NOTE 
	// When querying with struct, GORM will only query with non-zero fields, 
	// that means if your field’s value is 0, '', false or other zero values, it won’t be used to build query conditions

	// So pass accID as 0, when you want to fetch data irrespective of account_id
	err = r.DB.WithContext(ctx).Where(n).Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) GetNotificationsByStatus(ctx context.Context, st string, accID int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	status, err := models.StringToNotificationStatus(st)
	if err != nil {
		return nil, err
	}

	n := &models.Notification{
		AccountID: accID,
		Status: status,
	}

	err = r.DB.WithContext(ctx).Where(n).Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	return notifications, nil
}