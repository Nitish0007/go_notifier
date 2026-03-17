package notification

// import (
// 	"context"
	
// 	"gorm.io/gorm"
// )

// type NotificationBatchRepository struct {
// 	DB *gorm.DB
// }

// func NewNotificationBatchRepository(db *gorm.DB) *NotificationBatchRepository {
// 	return &NotificationBatchRepository{DB: db}
// }

// func (r *NotificationBatchRepository) Create(ctx context.Context, batch *NotificationBatch) error {
// 	return r.DB.WithContext(ctx).Create(&batch).Error
// }

// func (r *NotificationBatchRepository) GetByAccountID(ctx context.Context, accountID int) ([]*NotificationBatch, error) {
// 	var batches []*NotificationBatch
// 	err := r.DB.WithContext(ctx).Where("account_id = ?", accountID).Find(&batches).Error; if err != nil {
// 		return nil, err
// 	}
// 	return batches, nil
// }


// func (r *NotificationBatchRepository) GetByID(ctx context.Context, id string) (*NotificationBatch, error) {
// 	var batch NotificationBatch
// 	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&batch).Error; if err != nil {
// 		return nil, err
// 	}
// 	return &batch, nil
// }


// // func (r *NotificationBatchRepository) UpdateProgress(ctx context.Context, batchID string, successful, failed int) error {
// // 	return r.DB.WithContext(ctx).Model(&models.NotificationBatch{}).
// // 		Where("batch_id = ?", batchID)
// // 		Updates(map[string]interface{}{
// // 			"total_processed": gorm.Expr("total_processed + ?", successful+failed),
// // 			"successful_count": gorm.Expr("successful_count + ?", successful),
// // 			"failed_count": gorm.Expr("failed_count + ?", failed),
// // 			"status": "completed", // You can make this conditional
// // 	}).Error
// // }