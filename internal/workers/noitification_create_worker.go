package workers

import (
	"log"
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
	rbmq "github.com/rabbitmq/amqp091-go"
	"github.com/Nitish0007/go_notifier/internal/services"
	rabbitmq_utils "github.com/Nitish0007/go_notifier/utils/rabbitmq"
)

var (
	batchQueueName = "notification_batch"
	maxRetries = 5
	retryDelay = 1 * time.Minute
)

type NotificationBatchWorker struct {
	dbConn *gorm.DB
	rbmqConn *rbmq.Connection
	queue *rabbitmq_utils.Queue
	ctx context.Context
	blkNotificationService *services.BulkNotificationService
}

func NewNotificationBatchWorker(dbConn *gorm.DB, rbmqConn *rbmq.Connection, ctx context.Context, blkNotificationService *services.BulkNotificationService) *NotificationBatchWorker {
	queue, err := rabbitmq_utils.NewQueue(batchQueueName)
	if err != nil {
		log.Printf("Error creating queue: %v", err)
		return nil
	}
	
	return &NotificationBatchWorker{
		dbConn: dbConn,
		rbmqConn: rbmqConn,
		queue: queue,
		ctx: ctx,
		blkNotificationService: blkNotificationService,
	}
}

func (w *NotificationBatchWorker) Consume() {
	ch, err := rabbitmq_utils.CreateChannel(w.rbmqConn)
	if err != nil {
		log.Printf("Error creating channel: %v", err)
		return
	}
	defer ch.Close()
	
	msgs, err := ch.Consume(w.queue.Main.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("Error in consuming messages: %v", err)
		return
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var body map[string]any
			err := json.Unmarshal(msg.Body, &body)
			if err != nil {
				log.Printf("Error in unmarshalling body: %v", err)
				msg.Ack(false)
				w.queue.PushToDLQ(body)
				continue
			}


			batchID, ok := body["batch_id"].(string)
			if !ok {
				log.Printf("Batch ID is not a string: %v", body["batch_id"])
				w.queue.PushToDLQ(body)
				continue
			}

			for retryCount := 1; retryCount <= maxRetries; retryCount++ {
				err = w.blkNotificationService.ProcessBatch(w.ctx, batchID)
				if err != nil {
					log.Printf("Error in processing batch: %v", err)
					w.queue.PushToRetry(body)
					continue
				}
				time.Sleep(retryDelay)
			}
			w.queue.PushToDLQ(body)
			continue
		}
	}()
		
	<-forever
}
