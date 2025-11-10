// This worker will pick notifications that have status pending and does not have a key in redis, which means the notification is not yet scheduled to be sent
package workers

import (
	"context"
	"log"

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
	dbConn                 *gorm.DB
	rbmqConn               *rbmq.Connection
	queue                  *rabbitmq_utils.Queue
	ctx                    context.Context
	blkNotificationService *services.BulkNotificationService
}

func NewNotificationSchedulerWorker(db *gorm.DB, rbmqConn *rbmq.Connection, ctx context.Context, s *services.BulkNotificationService) *NotificationSchedulerWorker {
	q, err := rabbitmq_utils.NewQueue(deliveryQueueName)
	if err != nil {
		return nil
	}

	return &NotificationSchedulerWorker{
		dbConn:                 db,
		rbmqConn:               rbmqConn,
		queue:                  q,
		ctx:                    ctx,
		blkNotificationService: s,
	}
}

func (w *NotificationSchedulerWorker) Consume() {
	forever := make(chan bool)
	repo := repositories.NewNotificationRepository(w.dbConn)
	n := &models.Notification{
		AccountID: 0,
		Status: models.Pending,
	}
	
	notifications, err := repo.GetNotificationsByObject(w.ctx, n)
	if err != nil {
		log.Printf("Error Fetching notifications: %v", err)
		return
	}

	go func() {
		for _, n := range notifications {
			body, err := n.ToMap()
			if err != nil {
				log.Printf("Error converting to map: %v", err)
				continue
			}

			// push to queue
			if err := rabbitmq_utils.PushToQueue(w.queue.Main.Name, body); err != nil {
				log.Printf("Error pushing to queue: %v", err)
				continue
			}

			fieldsToUpdate := map[string]any {
				"status": models.Enqueued,
			}
			_, err = repo.UpdateNotification(w.ctx, fieldsToUpdate, n)
			if err != nil {
				log.Printf("Error updating notification: %v", err)
				continue
			}
		}
	}()
	<-forever
}
