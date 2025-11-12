package workers

import (
	"log"
	"context"
	"encoding/json"

	"gorm.io/gorm"
	rbmq "github.com/rabbitmq/amqp091-go"
	"github.com/Nitish0007/go_notifier/internal/services"
	rabbitmq_utils "github.com/Nitish0007/go_notifier/utils/rabbitmq"
)

var (
	notificationDeliveryQueueName = "notification_delivery"
)

type EmailWorker struct {
	dbConn *gorm.DB
	rbmqConn *rbmq.Connection
	ctx context.Context
	queue *rabbitmq_utils.Queue
	notificationService *services.NotificationService
}

func NewEmailWorker(dbConn *gorm.DB, rbmqConn *rbmq.Connection, ctx context.Context, queue *rabbitmq_utils.Queue, notificationService *services.NotificationService) *EmailWorker {
	q, err := rabbitmq_utils.NewQueue(notificationDeliveryQueueName)
	if err != nil {
		return nil
	}
	return &EmailWorker{
		dbConn: dbConn,
		rbmqConn: rbmqConn,
		ctx: ctx,
		queue: q,
		notificationService: notificationService,
	}
}


func (w *EmailWorker) Consume() {
	forever := make(chan bool)
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

			notificationID, ok := body["notificationID"].(string)
			if !ok {
				log.Printf("Notification ID is not a string: %v", body["notificationID"])
				w.queue.PushToDLQ(body)
				continue
			}

			accountID, ok := body["accountID"].(int)
			if !ok {
				log.Printf("Account ID is not an integer: %v", body["accountID"])
				w.queue.PushToDLQ(body)
				continue
			}


			err = w.notificationService.SendNotification(w.ctx, notificationID, accountID)
			if err != nil {
				log.Printf("Error in sending notification: %v", err)
				w.queue.PushToRetry(body)
				continue
			}

			msg.Ack(false)
			continue
		}
	}()
	<-forever
}