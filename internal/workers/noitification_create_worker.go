package workers

import (
	"encoding/json"
	"log"
	"context"

	"github.com/Nitish0007/go_notifier/internal/services"
	"github.com/Nitish0007/go_notifier/utils"
	rbmq "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

const (
	NOTIFICATION_CREATE_QUEUE = "notification_batch"
)

type NotificationBatchWorker struct {
	dbConn *gorm.DB
	rbmqConn *rbmq.Connection
	channel *rbmq.Channel
	queue *rbmq.Queue
	ctx context.Context
	blkNotificationService *services.BulkNotificationService
}

func NewNotificationBatchWorker(dbConn *gorm.DB, rbmqConn *rbmq.Connection, ctx context.Context, blkNotificationService *services.BulkNotificationService) *NotificationBatchWorker {
	channel, err := utils.CreateChannel(rbmqConn)
	if err != nil {
		log.Printf("Error creating channel: %v", err)
		return nil
	}
	defer channel.Close()

	queue, err := utils.CreateQueue(channel, NOTIFICATION_CREATE_QUEUE)
	if err != nil {
		log.Printf("Error creating queue: %v", err)
		return nil
	}

	return &NotificationBatchWorker{
		dbConn: dbConn,
		rbmqConn: rbmqConn,
		channel: channel,
		queue: queue,
		ctx: ctx,
		blkNotificationService: blkNotificationService,
	}
}

func (w *NotificationBatchWorker) Consume() {
	msgs, err := w.channel.Consume(w.queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("Error in worker: %v", err)
		return
	}

	log.Println("[*] Waiting for Job in queue. Press Ctrl+C to exit")

	forever := make(chan bool)
	go func() {
		log.Println("Goroutine started, waiting for messages...")
		for d := range msgs {
			log.Printf("======>> message : %v", d)
			var body map[string]any
			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				log.Printf("Error in unmarshalling body: %v", err)
				continue
			}

			err = w.blkNotificationService.ProcessBatch(w.ctx, body["batch_id"].(string))
			if err != nil {
				log.Printf("Error in processing batch: %v", err)
				continue
			}
			log.Printf("Batch processed successfully: %v", body["batch_id"])
			d.Ack(false)
		}
	}()

	<-forever
}
