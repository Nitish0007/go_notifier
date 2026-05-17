package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Nitish0007/go_notifier/internal/app/delivery"
	"github.com/Nitish0007/go_notifier/internal/common/mq"
	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
)

const emailDeliveryBase = "email_delivery"

type EmailWorker struct {
	mqClient          mq.MQClient
	ctx               context.Context
	campaignDeliverer *delivery.CampaignDeliverer
}

func NewEmailWorker(
	mqClient mq.MQClient,
	ctx context.Context,
	campaignDeliverer *delivery.CampaignDeliverer,
) *EmailWorker {
	return &EmailWorker{
		mqClient:          mqClient,
		ctx:               ctx,
		campaignDeliverer: campaignDeliverer,
	}
}

func (w *EmailWorker) RetryCount() int            { return 0 }
func (w *EmailWorker) MaxRetries() int            { return 5 }
func (w *EmailWorker) QueueName() string          { return emailDeliveryBase }
func (w *EmailWorker) RetryDelay() time.Duration  { return 1 * time.Minute }

func (w *EmailWorker) consumePolicy() *mq.ConsumePolicy {
	return &mq.ConsumePolicy{
		MaxRetries: w.MaxRetries(),
		BaseDelay:  w.RetryDelay(),
		MainQueue:  emailDeliveryBase,
		RetryQueue: emailDeliveryBase + "_retry",
		DLQQueue:   emailDeliveryBase + "_dlq",
	}
}

func (w *EmailWorker) Run() {
	policy := w.consumePolicy()
	retryQueue := emailDeliveryBase + "_retry"

	go func() {
		if err := sharedhelper.ConsumeFromMQ(w.ctx, w.mqClient, emailDeliveryBase, policy, w.handleDelivery); err != nil {
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

	notificationID, accountID, err := parseNotificationPayload(payload)
	if err != nil {
		return err
	}

	if err := w.campaignDeliverer.Deliver(ctx, notificationID, accountID); err != nil {
		return err
	}

	log.Printf("[email worker] delivered notification_id=%d account_id=%d", notificationID, accountID)
	return nil
}

func parseNotificationPayload(payload map[string]any) (notificationID int64, accountID int64, err error) {
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
		accountID = int64(v)
	case int:
		accountID = int64(v)
	case int64:
		accountID = v
	default:
		return 0, 0, fmt.Errorf("account id missing or invalid type %T", rawAID)
	}
	return notificationID, accountID, nil
}
