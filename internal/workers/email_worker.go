package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/Nitish0007/go_notifier/internal/common/mq"
	"github.com/Nitish0007/go_notifier/internal/features/configuration"
	"github.com/Nitish0007/go_notifier/internal/features/emailnotification"
	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
)

const notificationDeliveryBase = "notification_delivery"

type EmailWorker struct {
	dbConn              *gorm.DB
	mqClient            mq.MQClient
	ctx                 context.Context
	notificationService *emailnotification.EmailNotificationService
	configurationRepo   *configuration.ConfigurationRepository
}

func NewEmailWorker(
	dbConn *gorm.DB,
	mqClient mq.MQClient,
	ctx context.Context,
	notificationService *emailnotification.EmailNotificationService,
) *EmailWorker {
	return &EmailWorker{
		dbConn:              dbConn,
		mqClient:            mqClient,
		ctx:                 ctx,
		notificationService: notificationService,
		configurationRepo:   configuration.NewConfigurationRepository(dbConn),
	}
}

func (w *EmailWorker) RetryCount() int            { return 0 }
func (w *EmailWorker) MaxRetries() int            { return 5 }
func (w *EmailWorker) QueueName() string         { return notificationDeliveryBase }
func (w *EmailWorker) RetryDelay() time.Duration { return 1 * time.Minute }

func (w *EmailWorker) consumePolicy() *mq.ConsumePolicy {
	return &mq.ConsumePolicy{
		MaxRetries: w.MaxRetries(),
		BaseDelay:  w.RetryDelay(),
		MainQueue:  notificationDeliveryBase,
		RetryQueue: notificationDeliveryBase + "_retry",
		DLQQueue:   notificationDeliveryBase + "_dlq",
	}
}

// Run starts consumers on the main and retry queues; retry/DLQ/backoff are handled inside mq.Consume when policy is set.
func (w *EmailWorker) Run() {
	policy := w.consumePolicy()
	retryQueue := notificationDeliveryBase + "_retry"

	go func() {
		if err := sharedhelper.ConsumeFromMQ(w.ctx, w.mqClient, notificationDeliveryBase, policy, w.handleDelivery); err != nil {
			log.Printf("[email worker] main queue consumer exited: %v", err)
		}
	}()
	go func() {
		if err := sharedhelper.ConsumeFromMQ(w.ctx, w.mqClient, retryQueue, policy, w.handleDelivery); err != nil {
			log.Printf("[email worker] retry queue consumer exited: %v", err)
		}
	}()
}

func (w *EmailWorker) handleDelivery(ctx context.Context, body []byte) error {
	msg, err := sharedhelper.Decode(body)
	if err != nil {
		return fmt.Errorf("decode job: %w", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("decode payload: %w", err)
	}

	nid, accountID, err := parseNotificationPayload(payload)
	if err != nil {
		return err
	}

	configFilter := map[string]any{
		"account_id":            accountID,
		"config_type":           configuration.SMTPConfig,
		"default_configuration": true,
	}
	config, err := w.configurationRepo.GetByFields(ctx, configFilter)
	if err != nil {
		return fmt.Errorf("get configuration: %w", err)
	}
	if config == nil {
		return fmt.Errorf("configuration not found for account %d", accountID)
	}
	if config.ConfigType != configuration.SMTPConfig {
		return fmt.Errorf("configuration is not SMTP for account %d", accountID)
	}

	smtpConfig := &configuration.SMTPConfiguration{}
	jsonData, err := json.Marshal(config.Settings)
	if err != nil {
		return fmt.Errorf("marshal smtp settings: %w", err)
	}
	if err := json.Unmarshal(jsonData, smtpConfig); err != nil {
		return fmt.Errorf("unmarshal smtp settings: %w", err)
	}

	_ = nid
	_ = smtpConfig
	// Wire SendNotification when available; return nil simulates success for now.
	// if err := w.notificationService.SendNotification(...); err != nil { return err }
	return nil
}

func parseNotificationPayload(payload map[string]any) (notificationID int64, accountID int, err error) {
	rawNID, ok := payload["notification_id"]
	if !ok {
		rawNID = payload["notificationID"]
	}
	switch v := rawNID.(type) {
	case float64:
		notificationID = int64(v)
	case int64:
		notificationID = v
	default:
		return 0, 0, fmt.Errorf("notification id missing or invalid type %T", rawNID)
	}

	rawAID, ok := payload["account_id"]
	if !ok {
		rawAID = payload["accountID"]
	}
	switch v := rawAID.(type) {
	case float64:
		accountID = int(v)
	case int:
		accountID = v
	case int64:
		accountID = int(v)
	default:
		return 0, 0, fmt.Errorf("account id missing or invalid type %T", rawAID)
	}
	return notificationID, accountID, nil
}
