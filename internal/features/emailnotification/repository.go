package emailnotification

import (
	"fmt"
	"log"
	"context"

	"gorm.io/gorm"
)

type EmailNotificationRepository struct {
	DB *gorm.DB
}

func NewEmailNotificationRepository(conn *gorm.DB) *EmailNotificationRepository {
	return &EmailNotificationRepository{
		DB: conn,
	}
}

func (r *EmailNotificationRepository) Create(ctx context.Context, n *EmailNotification) error {
	return r.DB.WithContext(ctx).Create(n).Error
}

func (r *EmailNotificationRepository) Index(ctx context.Context, accID int) ([]*EmailNotification, error) {
	var notifications []*EmailNotification
	err := r.DB.WithContext(ctx).Where("account_id = ?", accID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *EmailNotificationRepository) GetByID(ctx context.Context, id string, accID int) (*EmailNotification, error) {
	var notification EmailNotification
	err := r.DB.WithContext(ctx).Where("id = ? AND account_id = ?", id, accID).First(&notification).Error
	return &notification, err
}

func (r *EmailNotificationRepository) GetNotificationsByStatus(ctx context.Context, st EmailNotificationStatus, accID int) ([]*EmailNotification, error) {
	var notifications []*EmailNotification

	n := &EmailNotification{
		AccountID: accID,
		Status:    st,
	}

	err := r.DB.WithContext(ctx).Where(n).Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *EmailNotificationRepository) GetNotificationsByObject(ctx context.Context, filters map[string]any, limit int) ([]*EmailNotification, error) {
	var notifications []*EmailNotification
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

func (r *EmailNotificationRepository) UpdateNotification(ctx context.Context, fieldsToUpdate map[string]any, nObj *EmailNotification) (*EmailNotification, error) {
	var udpatedNotification EmailNotification
	result := r.DB.WithContext(ctx).Model(&EmailNotification{}).Where("id = ? AND account_id = ?", nObj.ID, nObj.AccountID).Updates(fieldsToUpdate)

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
