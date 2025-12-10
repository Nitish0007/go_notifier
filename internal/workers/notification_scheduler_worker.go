// This worker(polling worker) will pick notifications that have status pending and does not have a key in redis, which means the notification is not yet scheduled to be sent. This will run every minute and enqueue 500 notifications at a time to avoid overwhelming the queue.
package workers

import (
	"context"
	"log"
	"time"

	"github.com/Nitish0007/go_notifier/internal/models"
	"github.com/Nitish0007/go_notifier/internal/repositories"
	"github.com/Nitish0007/go_notifier/internal/services"
	rabbitmq_utils "github.com/Nitish0007/go_notifier/utils/rabbitmq"
	rbmq "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

var (
	deliveryQueueName = "notification_delivery"
)

type NotificationSchedulerWorker struct {
	dbConn   *gorm.DB
	rbmqConn *rbmq.Connection
	queue    *rabbitmq_utils.Queue
	ctx      context.Context
}

func NewNotificationSchedulerWorker(db *gorm.DB, rbmqConn *rbmq.Connection, ctx context.Context, s *services.BulkNotificationService) *NotificationSchedulerWorker {
	q, err := rabbitmq_utils.NewQueue(deliveryQueueName)
	if err != nil {
		return nil
	}

	return &NotificationSchedulerWorker{
		dbConn:   db,
		rbmqConn: rbmqConn,
		queue:    q,
		ctx:      ctx,
	}
}

func (w *NotificationSchedulerWorker) Consume() {
	forever := make(chan bool)
	repo := repositories.NewNotificationRepository(w.dbConn)
	filters := map[string]any{
		// NOTE: this is the status of the notifications that are to be scheduled
		// using Enqueued status for development purposes
		"status": models.Pending,
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	go func() {

		for {
			select {
			case <-ticker.C:
				// enqueue 500 notifications at a time to avoid overwhelming the queue
				notifications, err := repo.GetNotificationsByObject(w.ctx, filters, 500)
				if err != nil {
					log.Printf("Error Fetching notifications: %v", err)
					continue
				}
				
				for _, n := range notifications {
					// create job message
					jobMessage := rabbitmq_utils.NewJobMessage(map[string]any{"notificationID": n.ID, "accountID": n.AccountID})
					// push to queue
					if err := rabbitmq_utils.PushToQueue(w.queue.Main, jobMessage); err != nil {
						log.Printf("Error pushing to queue: %v", err)
						continue
					}
		
					fieldsToUpdate := map[string]any{
						"status": models.Pending,
					}
					_, err = repo.UpdateNotification(w.ctx, fieldsToUpdate, n)
					if err != nil {
						log.Printf("Error updating notification: %v", err)
						continue
					}
				}

			case <-w.ctx.Done():
				log.Printf("Context done, exiting notification scheduler worker...")
				return
			}
		}
	}()
	<-forever
}
