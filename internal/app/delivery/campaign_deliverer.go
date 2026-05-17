package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Nitish0007/go_notifier/internal/features/configuration"
	"github.com/Nitish0007/go_notifier/internal/features/content"
	"github.com/Nitish0007/go_notifier/internal/features/emailnotification"
	libnotifier "github.com/Nitish0007/go_notifier/internal/lib/notifier"
	"github.com/Nitish0007/go_notifier/internal/shared/dto"
	notifierif "github.com/Nitish0007/go_notifier/internal/shared/interfaces/notifier"
)

// CampaignDeliverer loads campaign data and sends via the notifier registry.
type CampaignDeliverer struct {
	notifications *emailnotification.EmailNotificationRepository
	content       *content.ContentRepository
	configs       *configuration.ConfigurationRepository
	registry      *libnotifier.Registry
}

func NewCampaignDeliverer(
	notifications *emailnotification.EmailNotificationRepository,
	content *content.ContentRepository,
	configs *configuration.ConfigurationRepository,
	registry *libnotifier.Registry,
) *CampaignDeliverer {
	return &CampaignDeliverer{
		notifications: notifications,
		content:       content,
		configs:       configs,
		registry:      registry,
	}
}

func (d *CampaignDeliverer) Deliver(ctx context.Context, notificationID, accountID int64) error {
	n, err := d.notifications.GetByID(ctx, notificationID, accountID)
	if err != nil {
		return fmt.Errorf("get notification: %w", err)
	}
	if n.Status != emailnotification.Enqueued && n.Status != emailnotification.Scheduled {
		return fmt.Errorf("notification %d is not deliverable (status=%d)", notificationID, n.Status)
	}

	c, err := d.content.GetByID(ctx, accountID, n.ContentID)
	if err != nil {
		return fmt.Errorf("get content: %w", err)
	}

	recipients, err := d.notifications.ListCampaignRecipients(ctx, accountID, notificationID)
	if err != nil {
		return fmt.Errorf("list recipients: %w", err)
	}
	if len(recipients) == 0 {
		return fmt.Errorf("no active recipients for notification %d", notificationID)
	}

	providerCfg, err := d.loadSMTPConfig(ctx, accountID, n.FromEmail)
	if err != nil {
		return err
	}

	for _, rcpt := range recipients {
		req := toDeliveryRequest(n, rcpt, c.Body)
		if err := d.registry.Notify(ctx, notifierif.ChannelEmail, req, providerCfg); err != nil {
			return fmt.Errorf("send to %s: %w", rcpt.Email, err)
		}
	}

	now := time.Now().UTC()
	fields := map[string]any{
		"status":  emailnotification.Sent,
		"sent_at": now,
	}
	if _, err := d.notifications.UpdateNotification(ctx, fields, n); err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	return nil
}

func (d *CampaignDeliverer) loadSMTPConfig(ctx context.Context, accountID int64, fallbackFrom string) (*dto.SMTPConfiguration, error) {
	filter := map[string]any{
		"account_id":  accountID,
		"config_type": configuration.SMTPConfig,
		"is_default":  true,
	}
	cfg, err := d.configs.GetByFields(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get smtp configuration: %w", err)
	}

	var smtp configuration.SMTPConfiguration
	raw, err := json.Marshal(cfg.Settings)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &smtp); err != nil {
		return nil, err
	}

	from := smtp.From
	if from == "" {
		from = fallbackFrom
	}
	return &dto.SMTPConfiguration{
		Host:     smtp.Host,
		Port:     int64(smtp.Port),
		Username: smtp.Username,
		Password: smtp.Password,
		From:     from,
	}, nil
}

func toDeliveryRequest(n *emailnotification.EmailNotification, rcpt emailnotification.CampaignRecipient, body string) notifierif.DeliveryRequest {
	toName := rcpt.FirstName
	if rcpt.LastName != "" {
		if toName != "" {
			toName += " "
		}
		toName += rcpt.LastName
	}
	if toName == "" {
		toName = rcpt.Email
	}

	from := n.FromEmail
	replyTo := n.ReplyToEmail
	if replyTo == "" {
		replyTo = from
	}

	return notifierif.DeliveryRequest{
		AccountID: n.AccountID,
		Recipient: rcpt.Email,
		Subject:   n.Subject,
		Body:      body,
		HTMLBody:  body,
		Metadata: map[string]string{
			"from_email":     from,
			"from_name":      n.FromName,
			"to_name":        toName,
			"reply_to_email": replyTo,
			"reply_to_name":  n.FromName,
		},
	}
}
