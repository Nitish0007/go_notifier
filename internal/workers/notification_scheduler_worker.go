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
	// Recover from any panics to keep the worker running
	defer func() {
		if r := recover(); r != nil {
			log.Printf("--- [Scheduler Worker] PANIC recovered: %v\n", r)
			// Restart the worker after a panic
			log.Printf("--- [Scheduler Worker] Restarting worker after panic...\n")
			go w.Consume()
		}
	}()

	log.Printf("--- [Scheduler Worker] Starting notification scheduler worker...\n")

	pollInterval := 1 * time.Minute

	// Run immediately on start
	log.Printf("--- [Scheduler Worker] Running initial poll...\n")
	err := w.pollNotifications()
	if err != nil {
		log.Printf(">>>>> Error polling notifications on startup: %v", err)
	} else {
		log.Printf("--- [Scheduler Worker] Initial poll completed successfully\n")
	}

	// Use time.After instead of ticker to ensure we wait the full interval AFTER each poll
	// This prevents overlapping polls if pollNotifications() takes longer than the interval
	for {
		log.Printf("--- [Scheduler Worker] Waiting %v until next poll...\n", pollInterval)
		select {
		case <-time.After(pollInterval):
			log.Printf("--- [Scheduler Worker] Interval elapsed, starting poll...\n")
			err := w.pollNotifications()
			if err != nil {
				log.Printf(">>>>> Error polling notifications: %v", err)
			} else {
				log.Printf("--- [Scheduler Worker] Poll completed successfully\n")
			}
			// Loop continues to wait for next interval
		case <-w.ctx.Done():
			log.Printf("--- [Scheduler Worker]	Context done, exiting notification scheduler worker...\n")
			return
		}
	}
}

// PRIVATE FUNCTIONS //

func (w *NotificationSchedulerWorker) pollNotifications() error {
	// fetch notifications from db using polling and push them to queue every minute
	notifications, err := w.fetchNotifications()
	if err != nil {
		log.Printf(">>>>> Error fetching notifications: %v", err)
		return err
	}
	err = w.pushNotificationsToDeliveryQueue(notifications)
	if err != nil {
		log.Printf(">>>>> Error pushing notifications to delivery queue: %v", err)
		return err
	}

	log.Printf(">>>>> Notifications pushed to delivery queue successfully")
	return nil
}

func (w *NotificationSchedulerWorker) fetchNotifications() ([]*models.Notification, error) {
	repo := repositories.NewNotificationRepository(w.dbConn)
	filters := map[string]any{
		// NOTE: this is the status of the notifications that are to be scheduled
		// using Enqueued status for development purposes
		"status": models.Pending,
	}

	// enqueue 500 notifications at a time to avoid overwhelming the queue
	log.Printf("--- [Scheduler Worker]	Fetching notifications..........\n")
	notifications, err := repo.GetNotificationsByObject(w.ctx, filters, 500)
	if err != nil {
		log.Printf("--- [Scheduler Worker]	Error Fetching notifications: %v\n", err)
		return nil, err
	}

	log.Printf("--- [Scheduler Worker]	Found %d notifications to be scheduled\n", len(notifications))
	return notifications, nil
}

func (w *NotificationSchedulerWorker) pushNotificationsToDeliveryQueue(notifications []*models.Notification) error {
	repo := repositories.NewNotificationRepository(w.dbConn)
	for _, n := range notifications {
		// create job message
		payload := map[string]any{"notificationID": n.ID, "accountID": n.AccountID}
		jobMessage := rabbitmq_utils.NewJobMessage(map[string]any{"payload": payload})
		jobMetaData := rabbitmq_utils.NewJobMetadata(0, rabbitmq_utils.MAX_RETRIES, rabbitmq_utils.RETRY_DELAY)
		ctx, cancel := context.WithTimeout(w.ctx, 5*time.Second)
		defer cancel()
		err := rabbitmq_utils.StoreJobMetadata(ctx, jobMessage.GetJobID(), *jobMetaData)
		if err != nil {
			log.Printf("--- [Scheduler Worker]	Error storing job metadata: %v\n", err)
			continue
		}

		// push to queue
		if err := rabbitmq_utils.PushToQueue(w.queue.Main, jobMessage); err != nil {
			log.Printf("--- [Scheduler Worker]	Error pushing to queue: %v\n", err)
			continue
		}

		fieldsToUpdate := map[string]any{
			"status": models.Enqueued,
		}
		_, err = repo.UpdateNotification(w.ctx, fieldsToUpdate, n)
		if err != nil {
			log.Printf("--- [Scheduler Worker]	Error updating notification: %v\n", err)
			continue
		}
	}
	return nil
}
