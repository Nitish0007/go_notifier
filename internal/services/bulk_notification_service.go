package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/Nitish0007/go_notifier/internal/models"
	"github.com/Nitish0007/go_notifier/internal/repositories"
	"github.com/Nitish0007/go_notifier/utils"
)

type BulkNotificationService struct {
	notificationRepo *repositories.NotificationRepository
	notificationBatchRepo *repositories.NotificationBatchRepository
	notificationBatchErrorRepo *repositories.NotificationBatchErrorRepo
}

func NewBulkNotificationService(nr *repositories.NotificationRepository, nbr *repositories.NotificationBatchRepository, nber *repositories.NotificationBatchErrorRepo) *BulkNotificationService {
	return &BulkNotificationService{
		notificationRepo: nr,
		notificationBatchRepo: nbr,
		notificationBatchErrorRepo: nber,
	}
}

func (s *BulkNotificationService) CreateBulkNotifications(ctx context.Context, data map[string]any) (map[string]any, error){
	// create a batch
	notificationRaw, exists := data["notifications"]
	if !exists {
		return nil, errors.New("notifications object is required in payload")
	}

	notifications, err := json.Marshal(notificationRaw)  // converting to json format
	if err != nil {
		return nil, errors.New("unable to marshal notifications")
	}

	var notificationsList []map[string]any
	err = json.Unmarshal(notifications, &notificationsList)
	if err != nil {
		return nil, errors.New("unable to unmarshal notifications")
	}

	nChannel, err := models.StringToNotificationChannel(notificationsList[0]["channel"].(string))
	if err != nil {
		return nil, err
	}
	batch := &models.NotificationBatch{
		AccountID: utils.GetCurrentAccountID(ctx),
		Count: len(data),
		Channel: nChannel, // assuming all notifications have the same channel
		Status: int(models.Pending), // pending
		Payload: data,
	}
	err = s.notificationBatchRepo.Create(ctx, batch)
	if err != nil {
		return nil, err
	}
	// push batch to queue for processing
	body := map[string]any{
		"batch_id": batch.ID,
		"account_id": utils.GetCurrentAccountID(ctx),
	}
	err = utils.PushToQueue("notification_batch", body)
	if err != nil {
		return nil, errors.New("failed to push batch to queue")
	}

	// return success response
	return map[string]any{
		"batch_id": batch.ID,
		"message": "batch created successfully and will be processed asynchronously",
	}, nil
}

func (s *BulkNotificationService) ProcessBatch(ctx context.Context, batchID string) {
	batch, err := s.notificationBatchRepo.GetByID(ctx, batchID)
	if err != nil {
		log.Printf("ERROR!: %v", err)
		return
	}

	data, ok := batch.Payload["notifications"].([]map[string]any)
	if !ok {
		log.Printf("ERROR!: bad format of payload: %v", batch.Payload)
		return
	}
	err = validateBulkNotificationPayload(data, batch.ID)
	if err != nil {
		log.Printf("ERROR!: %v", err)
		return
	}
}

// Private method
func validateBulkNotificationPayload(payload []map[string]any, batchID string) error {
	payloadCollection := make(chan map[string]any)   // contain all payload
	// validPayloads := make(chan bool)   // valid chunks of payloads
	// invalidPayloads := make(chan bool) // chunks of invalid payloads

	var workers int
	wg := sync.WaitGroup{}

	// set the number of workers based on the payload size
	switch payloadSize := len(payload); {
	case payloadSize < 5:
		workers = 1
	case payloadSize > 5 && payloadSize < 60:
		workers = 3
	default:
		workers = 5
	}

	// worker pool pattern
	// Initializing go routines(workers)
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range payloadCollection {
				_, err := utils.ValidateNotificationPayload(p)
				if err != nil {
					// create notification batch error
					// 	notificationBatchError := &models.NotificationBatchError{
					// 		BatchID: batchID,
					// 		BatchIndex: p["batch_index"].(int),
					// 		AccountID: utils.GetCurrentAccountID(ctx),
					// 		ErrorMessage: err.Error(),
					// 		Payload: p,
					// 	}
					// 	err = s.notificationBatchErrorRepo.Create(ctx, notificationBatchError)
					// 	if err != nil {
					// 		log.Printf("ERROR!: %v", err)
					// 		return
					// 	}
					// } else {
					// 	// create notification
					// 	_, err := s.notificationRepo.Create(context.Background(), &models.Notification{
					// 		BatchID: &batchID,
					// 		Channel: vp["channel"].(int),
					// 		Recipient: vp["recipient"].(string),
					// 		Subject: vp["subject"].(string),
					// 		Body: vp["body"].(string),
					// 		HtmlBody: vp["html_body"].(string),
					// 		Status: 0,
					// 		Metadata: vp["metadata"].(map[string]any),
					// 		ErrorMessage: nil,
					// 		JobID: nil,
					// 		SendAt: vp["send_at"].(*time.Time),
					// 		SentAt: nil,
					// 		CreatedAt: time.Now(),
					// 		UpdatedAt: time.Now(),
					// 	})
					// 	if err != nil {
					// 		log.Printf("ERROR!: %v", err)
					// 		return
					// 	}
				}
			}
		}()
	}
	return nil
}