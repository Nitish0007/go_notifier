package emailnotificationlist

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
)

type EmailNotificationListRepository struct {
	db *gorm.DB
}

func NewEmailNotificationListRepository(db *gorm.DB) *EmailNotificationListRepository {
	return &EmailNotificationListRepository{db: db}
}

// EnsureLinked replaces all email_notification_lists rows for this notification with listIDs.
// Existing links not in listIDs are removed; listIDs are inserted. Pass an empty slice to clear all links.
func (r *EmailNotificationListRepository) EnsureLinked(ctx context.Context, tx *gorm.DB, accountID int64, listIDs []int64, notificationID int64) error {
	if tx == nil {
		tx = r.db
	}
	ids := sharedhelper.GetUniques(listIDs)

	if len(ids) > 0 {
		var count int64
		if err := tx.WithContext(ctx).Table("lists").Where("account_id = ? AND id IN ?", accountID, ids).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(ids)) {
			return errors.New("invalid list id/ids passed")
		}
	}

	if err := tx.WithContext(ctx).
		Where("notification_id = ? AND account_id = ?", notificationID, accountID).
		Delete(&EmailNotificationList{}).Error; err != nil {
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	rows := make([]EmailNotificationList, 0, len(ids))
	for _, lid := range ids {
		rows = append(rows, *NewEmailNotificationList(accountID, lid, notificationID))
	}
	return tx.WithContext(ctx).CreateInBatches(rows, 100).Error
}