package workers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Nitish0007/go_notifier/internal/models"
	"github.com/Nitish0007/go_notifier/internal/repositories"
	"github.com/Nitish0007/go_notifier/internal/services"
	rabbitmq_utils "github.com/Nitish0007/go_notifier/utils/rabbitmq"
	rbmq "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

var (
	notificationDeliveryQueueName = "notification_delivery"
)

type EmailWorker struct {
	dbConn              *gorm.DB
	rbmqConn            *rbmq.Connection
	ctx                 context.Context
	queue               *rabbitmq_utils.Queue
	notificationService *services.NotificationService
	configurationRepo   *repositories.ConfigurationRepository
}

func NewEmailWorker(dbConn *gorm.DB, rbmqConn *rbmq.Connection, ctx context.Context, notificationService *services.NotificationService) *EmailWorker {
	q, err := rabbitmq_utils.NewQueue(notificationDeliveryQueueName)
	if err != nil {
		return nil
	}
	configurationRepo := repositories.NewConfigurationRepository(dbConn)
	return &EmailWorker{
		dbConn:              dbConn,
		rbmqConn:            rbmqConn,
		ctx:                 ctx,
		queue:               q,
		notificationService: notificationService,
		configurationRepo:   configurationRepo,
	}
}

func (w *EmailWorker) Consume() {
	log.Printf(">>>>>>>>>>>>>>>>> Consuming email notifications\n")
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

			log.Printf("========================> Message body without unmarshalling: %v", string(msg.Body))
			jobMsg := rabbitmq_utils.NewJobMessage(map[string]any{})
			err := jobMsg.FromJSON(msg.Body)
			if err != nil {
				log.Printf("Error in unmarshalling body: %v", err)
				w.queue.PushToDLQ(jobMsg)
				msg.Ack(false)
				continue
			}

			body := jobMsg.GetPayload()
			log.Printf("========================> Message body: %v", body)
			notificationID, ok := body["notificationID"].(string)
			if !ok {
				log.Printf("Notification ID not found in message body %v", body)
				w.queue.PushToDLQ(jobMsg)
				msg.Ack(false)
				continue
			}
			
			aid, exists := body["accountID"]
			if !exists {
				log.Printf("Account ID not found in message body %v", body)
				w.queue.PushToDLQ(jobMsg)
				msg.Ack(false)
				continue
			}
			
			accountIDFloat64, ok := aid.(float64)
			if !ok {
				log.Printf("Account ID is not a number: %v", aid)
				w.queue.PushToDLQ(jobMsg)
				continue
			}
			
			accountID := int(accountIDFloat64)
			configFilter := map[string]any{
				"account_id":            accountID,
				"config_type":           string(models.SMTPConfig),
				"default_configuration": true,
			}
			config, err := w.configurationRepo.GetByFields(w.ctx, configFilter)
			if err != nil {
				log.Printf("Error in getting configuration: %v", err)
				w.queue.PushToDLQ(jobMsg)
				msg.Ack(false)
				continue
			}
			if config == nil {
				log.Printf("Configuration not found for account ID: %v", accountID)
				w.queue.PushToDLQ(jobMsg)
				msg.Ack(false)
				continue
			}
			if config.ConfigType != string(models.SMTPConfig) {
				log.Printf("Configuration is not an email configuration: %v", config.ConfigType)
				w.queue.PushToDLQ(jobMsg)
				msg.Ack(false)
				continue
			}
			smtpConfig := &models.SMTPConfiguration{}
			jsonData, err := json.Marshal(config.ConfigurationData)
			if err != nil {
				log.Printf("Error in marshalling SMTP configuration: %v", err)
				w.queue.PushToDLQ(jobMsg)
				continue
			}
			err = json.Unmarshal(jsonData, smtpConfig)
			if err != nil {
				log.Printf("Error in unmarshalling SMTP configuration: %v", err)
				w.queue.PushToDLQ(jobMsg)
				msg.Ack(false)
				continue
			}
			err = w.notificationService.SendNotification(w.ctx, notificationID, accountID, smtpConfig)
			if err != nil {
				log.Printf("Error in sending notification: %v", err)
				w.queue.PushToRetry(jobMsg)
				msg.Ack(false)
				continue
			}
			msg.Ack(false)
		}
	}()
	<-forever
	log.Printf(">>>>>>>>>>>>>>>>> Consumed email notifications\n\n")
}

func (w *EmailWorker) ConsumeRetry() {
	log.Printf(">>>>>>>>>>>>>>>>> Consuming retry email notifications\n")
	forever := make(chan bool)
	ch, err := rabbitmq_utils.CreateChannel(w.rbmqConn)
	if err != nil {
		log.Printf("Error creating channel: %v", err)
		return
	}
	defer ch.Close()

	msgs, err := ch.Consume(w.queue.Retry.Name, "", true, false, false, false, nil)
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
				w.queue.PushToDLQ(rabbitmq_utils.NewJobMessage(body))
				continue
			}
		}
	}()
	<-forever
	log.Printf(">>>>>>>>>>>>>>>>> Consumed retry email notifications\n\n")
}
