package repositories

import (
	"context"
	
	"github.com/Nitish0007/go_notifier/internal/models"
	"gorm.io/gorm"
)

type NotificationBatchRepository struct {
	DB *gorm.DB
}

func NewNotificationBatchRepository(db *gorm.DB) *NotificationBatchRepository {
	return &NotificationBatchRepository{DB: db}
}

func (r *NotificationBatchRepository) Create(ctx context.Context, batch *models.NotificationBatch) error {
	return r.DB.WithContext(ctx).Create(batch).Error
}

func (r *NotificationBatchRepository) GetByBatchID(ctx context.Context, batchID string) (*models.NotificationBatch, error) {
	var batch models.NotificationBatch
	err := r.DB.WithContext(ctx).Where("batch_id = ?", batchID).First(&batch).Error
	return &batch, err
}

// func (r *NotificationBatchRepository) UpdateProgress(ctx context.Context, batchID string, successful, failed int) error {
// 	return r.DB.WithContext(ctx).Model(&models.NotificationBatch{}).
// 		Where("batch_id = ?", batchID)
// 		Updates(map[string]interface{}{
// 			"total_processed": gorm.Expr("total_processed + ?", successful+failed),
// 			"successful_count": gorm.Expr("successful_count + ?", successful),
// 			"failed_count": gorm.Expr("failed_count + ?", failed),
// 			"status": "completed", // You can make this conditional
// 	}).Error
// }