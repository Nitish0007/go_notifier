package services

import (
	"context"

	"github.com/Nitish0007/go_notifier/internal/repositories"
)

type BulkNotificationService struct {
	notificationRepo *repositories.NotificationRepository
	notificationBatchRepo *repositories.NotificationBatchRepository
	notificationBatchErrorRepo *repositories.NotificationBatchErrorRepo
}

func NewBulkNotificationService(notificationRepo *repositories.NotificationRepository) *BulkNotificationService {
	return &BulkNotificationService{
		notificationRepo: notificationRepo,
	}
}

func (s *BulkNotificationService) CreateBulkNotifications(ctx context.Context, data []map[string]any) (map[string]any, error){
	// create a batch
	// batch := s.notificationBatchRepo.Create(ctx)
	
	// go func(){

	// 	validPayloads, invalidPayloads := utils.ValidateBulkNotificationPayload(data)
	// 	if len(invalidPayloads) != 0 {
	// 		// handle invalid payloads and return errors for them
	// 		// if valid payloads are present then they should be processed and notifications should be sent
	
	// 	}
		
	// 	// create for validPayloads
	// 	// assign batch id to all payloads
	// 	for _, vp := range validPayloads {
	// 		vp["batch_id"] = batchId
			
	// 	}
	// }()
	return nil, nil
}