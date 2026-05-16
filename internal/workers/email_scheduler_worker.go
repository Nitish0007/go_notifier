// This worker(polling worker) will pick notifications that have status pending and does not have a key in redis, which means the notification is not yet scheduled to be sent. This will run every minute and enqueue 500 notifications at a time to avoid overwhelming the queue.
package workers

import (
	"log"
	"time"
	"context"
	"gorm.io/gorm"
	"github.com/Nitish0007/go_notifier/internal/common/mq"
	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
	"github.com/Nitish0007/go_notifier/internal/features/emailnotification"
	"github.com/Nitish0007/go_notifier/internal/features/emailnotificationlist"
)

type EmailSchedulerWorker struct {
	dbConn   *gorm.DB
	mqClient mq.MQClient
	ctx      context.Context
}

func NewEmailSchedulerWorker(ctx context.Context, db *gorm.DB, mqClient mq.MQClient) *EmailSchedulerWorker {
	return &EmailSchedulerWorker{
		dbConn:   db,
		mqClient: mqClient,
		ctx:      ctx,
	}
}

func (w *EmailSchedulerWorker) RetryCount() int { return 0 }
func (w *EmailSchedulerWorker) MaxRetries() int { return 0 }
func (w *EmailSchedulerWorker) QueueName() string { return "email_scheduler_queue" }
func (w *EmailSchedulerWorker) RetryDelay() time.Duration {	return 1 * time.Minute }
func (w *EmailSchedulerWorker) pollInterval() time.Duration { return 1 * time.Minute }


func (w *EmailSchedulerWorker) Run() {
	// Recover from any panics to keep the worker running
	defer func() {
		if r := recover(); r != nil {
			log.Printf("--- [Email Scheduler Worker] PANIC recovered: %v\n", r)
			// Restart the worker after a panic
			log.Printf("--- [Email Scheduler Worker] Restarting worker after panic...\n")
			go w.Run()
		}
	}()

	log.Printf("--- [Email Scheduler Worker] Starting notification scheduler worker...\n")

	// Run immediately on start
	log.Printf("--- [Email Scheduler Worker] Running initial poll...\n")
	err := w.pollNotifications()
	if err != nil {
		log.Printf(">>>>> Error polling notifications on startup: %v", err)
	} else {
		log.Printf("--- [Scheduler Worker] Initial poll completed successfully\n")
	}

	// Use time.After instead of ticker to ensure we wait the full interval AFTER each poll
	// This prevents overlapping polls if pollNotifications() takes longer than the interval
	for {
		log.Printf("--- [Email Scheduler Worker] Waiting %v until next poll...\n", w.pollInterval())
		select {
		case <-time.After(w.pollInterval()):
			log.Printf("--- [Email Scheduler Worker] Interval elapsed, starting poll...\n")
			err := w.pollNotifications()
			if err != nil {
				log.Printf(">>>>> Error polling email notifications: %v", err)
			} else {
				log.Printf("--- [Email Scheduler Worker] Poll completed successfully\n")
			}
			// Loop continues to wait for next interval
		case <-w.ctx.Done():
			log.Printf("--- [Email Scheduler Worker]	Context done, exiting email scheduler worker...\n")
			return
		}
	}
}

// PRIVATE FUNCTIONS //

func (w *EmailSchedulerWorker) pollNotifications() error {
	// fetch notifications from db using polling and push them to queue every minute
	notifications, err := w.fetchNotifications()
	if err != nil {
		log.Printf(">>>>> Error fetching notifications: %v", err)
		return err
	}
	
	if len(notifications) == 0 {
		log.Printf("--- [Email Scheduler Worker]	No notifications to be scheduled\n")
		return nil
	}
	err = w.pushNotificationsToDeliveryQueue(notifications)
	if err != nil {
		log.Printf(">>>>> Error pushing notifications to delivery queue: %v", err)
		return err
	}

	log.Printf(">>>>> Notifications pushed to delivery queue successfully")
	return nil
}

func (w *EmailSchedulerWorker) fetchNotifications() ([]*emailnotification.EmailNotification, error) {
	listLinks := emailnotificationlist.NewEmailNotificationListRepository(w.dbConn)
	repo := emailnotification.NewEmailNotificationRepository(w.dbConn, listLinks)
	query := repo.DB.WithContext(w.ctx).Where(
		"status = ? AND notification_type = ? AND send_at IS NOT NULL AND send_at <= ?",
		emailnotification.Scheduled,
		emailnotification.Campaign,
		time.Now(),
	).Limit(500)

	// enqueue 500 notifications at a time to avoid overwhelming the queue
	log.Printf("--- [Email Scheduler Worker]	Fetching notifications..........\n")
	var notifications []*emailnotification.EmailNotification
	err := query.Find(&notifications).Error
	if err != nil {
		log.Printf("--- [Email Scheduler Worker]	Error Fetching notifications: %v\n", err)
		return nil, err
	}

	log.Printf("--- [Email Scheduler Worker]	Found %d email notifications to be scheduled\n", len(notifications))
	return notifications, nil
}

func (w *EmailSchedulerWorker) pushNotificationsToDeliveryQueue(notifications []*emailnotification.EmailNotification) error {
	listLinks := emailnotificationlist.NewEmailNotificationListRepository(w.dbConn)
	repo := emailnotification.NewEmailNotificationRepository(w.dbConn, listLinks)
	for _, n := range notifications {
		// create job metadata
		metaData := sharedhelper.JobMetadata{
			RetryCount: w.RetryCount(),
			MaxRetries: w.MaxRetries(),
			RetryDelay: w.RetryDelay(),
		}

		// create job message
		payload := map[string]any{"notification_id": n.ID, "account_id": n.AccountID}
		jobMessage, err := sharedhelper.NewMQMessage(payload, &metaData)
		if err != nil {
			log.Printf("--- [Email Scheduler Worker]	Error creating job message: %v\n", err)
			continue
		}

		// push to queue
		if err := sharedhelper.PublishToMQ(w.ctx, w.mqClient, w.QueueName(), jobMessage); err != nil {
			log.Printf("--- [Email Scheduler Worker]	Error publishing to queue: %v\n", err)
			continue
		}

		fieldsToUpdate := map[string]any{
			"status": emailnotification.Enqueued,
		}

		_, err = repo.UpdateNotification(w.ctx, fieldsToUpdate, n)
		if err != nil {
			log.Printf("--- [Email Scheduler Worker]	Error updating notification: %v\n", err)
			continue
		}
	}
	return nil
}