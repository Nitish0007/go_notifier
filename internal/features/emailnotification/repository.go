package emailnotification

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Nitish0007/go_notifier/internal/features/emailnotificationlist"
	"gorm.io/gorm"
)

type EmailNotificationRepository struct {
	DB        *gorm.DB
	emailNotificationListRepo *emailnotificationlist.EmailNotificationListRepository
}

func NewEmailNotificationRepository(conn *gorm.DB, enlr *emailnotificationlist.EmailNotificationListRepository) *EmailNotificationRepository {
	return &EmailNotificationRepository{
		DB:        conn,
		emailNotificationListRepo: enlr,
	}
}

func (r *EmailNotificationRepository) Create(ctx context.Context, n *EmailNotification) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return r.createWithTx(ctx, tx, n)
	})
}

// CreateCampaignWithList creates an email notification and links it to a list in one transaction.
func (r *EmailNotificationRepository) CreateCampaignWithList(ctx context.Context, n *EmailNotification, listIDs []int64) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.createWithTx(ctx, tx, n); err != nil {
			return err
		}
		return r.emailNotificationListRepo.EnsureLinked(ctx, tx, n.AccountID, listIDs, n.ID)
	})
}

func (r *EmailNotificationRepository) createWithTx(ctx context.Context, tx *gorm.DB, n *EmailNotification) error {
	if err := tx.WithContext(ctx).Create(n).Error; err != nil {
		return errors.New("failed to create email notification: " + err.Error())
	}
	return nil
}

func (r *EmailNotificationRepository) Index(ctx context.Context, accID int64) ([]*EmailNotification, error) {
	var notifications []*EmailNotification
	err := r.DB.WithContext(ctx).Where("account_id = ?", accID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *EmailNotificationRepository) GetByID(ctx context.Context, id int64, accID int64) (*EmailNotification, error) {
	var notification EmailNotification
	err := r.DB.WithContext(ctx).Where("id = ? AND account_id = ?", id, accID).First(&notification).Error
	return &notification, err
}

func (r *EmailNotificationRepository) GetNotificationsByStatus(ctx context.Context, st EmailNotificationStatus, accID int64) ([]*EmailNotification, error) {
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

func (r *EmailNotificationRepository) ListCampaignRecipients(ctx context.Context, accountID, notificationID int64) ([]CampaignRecipient, error) {
	var rows []CampaignRecipient
	err := r.DB.WithContext(ctx).Raw(`
		SELECT DISTINCT ec.email,
			COALESCE(c.first_name, '') AS first_name,
			COALESCE(c.last_name, '') AS last_name
		FROM email_notification_lists enl
		INNER JOIN list_subscriptions ls
			ON ls.list_id = enl.list_id AND ls.account_id = enl.account_id AND ls.active = true
		INNER JOIN contacts c
			ON c.id = ls.contact_id AND c.account_id = ls.account_id
		INNER JOIN email_contacts ec
			ON ec.contact_id = c.id AND ec.account_id = c.account_id
		WHERE enl.notification_id = ? AND enl.account_id = ?
	`, notificationID, accountID).Scan(&rows).Error
	return rows, err
}
