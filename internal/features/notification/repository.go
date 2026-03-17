package notification

import (
	"fmt"
	"log"
	"context"

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

func (r *NotificationRepository) Create(ctx context.Context, n *Notification) error {
	return r.DB.WithContext(ctx).Create(n).Error
}

func (r *NotificationRepository) Index(ctx context.Context, accID int) ([]*Notification, error) {
	var notifications []*Notification
	err := r.DB.WithContext(ctx).Where("account_id = ?", accID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) GetByID(ctx context.Context, id string, accID int) (*Notification, error) {
	var notification Notification
	err := r.DB.WithContext(ctx).Where("id = ? AND account_id = ?", id, accID).First(&notification).Error
	return &notification, err
}

func (r *NotificationRepository) GetNotificationsByChannel(ctx context.Context, ch NotificationChannel, accID int) ([]*Notification, error) {
	var notifications []*Notification

	n := &Notification{
		AccountID: accID,
		Channel:   ch,
	}

	// NOTE
	// When querying with struct, GORM will only query with non-zero fields,
	// that means if your field’s value is 0, '', false or other zero values, it won’t be used to build query conditions

	// So pass accID as 0, when you want to fetch data irrespective of account_id
	err := r.DB.WithContext(ctx).Where(n).Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) GetNotificationsByStatus(ctx context.Context, st NotificationStatus, accID int) ([]*Notification, error) {
	var notifications []*Notification

	n := &Notification{
		AccountID: accID,
		Status:    st,
	}

	err := r.DB.WithContext(ctx).Where(n).Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *NotificationRepository) GetNotificationsByObject(ctx context.Context, filters map[string]any, limit int) ([]*Notification, error) {
	var notifications []*Notification
	query := r.DB.WithContext(ctx)
	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) UpdateNotification(ctx context.Context, fieldsToUpdate map[string]any, nObj *Notification) (*Notification, error) {
	var udpatedNotification Notification
	result := r.DB.WithContext(ctx).Model(&Notification{}).Where("id = ? AND account_id = ?", nObj.ID, nObj.AccountID).Updates(fieldsToUpdate)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>	No rows affected while updating notification")
		return nObj, nil
	}

	// fetch updated record
	err := r.DB.WithContext(ctx).Where("id = ? AND account_id = ?", nObj.ID, nObj.AccountID).First(&udpatedNotification).Error
	if err != nil {
		return nil, err
	}

	return &udpatedNotification, nil
}
