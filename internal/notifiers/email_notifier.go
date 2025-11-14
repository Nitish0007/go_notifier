package notifiers

import (
	"context"
	"errors"

	"github.com/Nitish0007/go_notifier/internal/models"
	"github.com/Nitish0007/go_notifier/internal/repositories"
	"github.com/Nitish0007/go_notifier/utils"
)

type EmailNotifier struct {
	notificationRepository *repositories.NotificationRepository
};

func NewEmailNotifier(r *repositories.NotificationRepository) *EmailNotifier {
	return &EmailNotifier{
		notificationRepository: r,
	}
}

func (n *EmailNotifier) Send(notification *models.Notification) error {
	
	return nil
}

func (n *EmailNotifier) CreateNotification(ctx context.Context, payload map[string]any) (*models.Notification, error) {
	// Initialize the notification model instance for given payload
	notification := &models.Notification{}
	aid := utils.GetCurrentAccountID(ctx)
	if aid == -1 {
		return nil, errors.New("account ID is required to create a notification")
	}
	notification.AccountID = aid
	emailChannel, exists := payload["channel"].(string)
	if !exists || emailChannel != "email" {
		return nil, errors.New("channel is required and must be 'email'")
	}

	nChannel, err := models.StringToNotificationChannel(emailChannel)
	if err != nil {
		return nil, err
	}

	// NOTE: THE PAYLOAD STRUCTURE MUST BE REFINED LATER, THIS IS JUST A PROTOTYPE
	notification.Channel = nChannel
	nStatus, err := models.StringToNotificationStatus("pending")
	if err != nil {
		return nil, err
	}
	notification.Status = nStatus
	notification.Recipient = payload["to"].(string)
	notification.Body = payload["body"].(string)
	notification.HtmlBody = payload["html_body"].(string)
	sendTime := payload["send_at"]
	if sendTime != nil {
		notification.SendAt = utils.ParseTime(sendTime.(string))
	} else {
		notification.SendAt = utils.GetCurrentTime()
	}
	
	sanitizedMetadata := make(map[string]any)	
	if metadata, exists := payload["metadata"].(map[string]any); exists {
		if len(metadata) == 0 {
			return nil, errors.New("metadata can't be empty")
		}
		if _, ok := metadata["from_email"]; !ok || metadata["from_email"] == "" {
			return nil, errors.New("from_email is required in metadata")
		}
		if _, ok := metadata["from_name"]; !ok || metadata["from_name"] == "" {
			return nil, errors.New("from_name is required in metadata")
		}
		if _, ok := metadata["to_name"]; !ok || metadata["to_name"] == "" {
			return nil, errors.New("to_name is required in metadata")
		}
		if _, ok := metadata["reply_to_email"]; !ok || metadata["reply_to_email"] == "" {
			return nil, errors.New("reply_to_email is required in metadata")
		}
		if _, ok := metadata["reply_to_name"]; !ok || metadata["reply_to_name"] == "" {
			return nil, errors.New("reply_to_name is required in metadata")
		}
		if _, ok := metadata["subject"]; !ok || metadata["subject"] == "" {
			return nil, errors.New("subject is required in metadata")
		}

		// keeping only required fields in metadata
		sanitizedMetadata["from_email"] = metadata["from_email"].(string)
		sanitizedMetadata["from_name"] = metadata["from_name"].(string)
		sanitizedMetadata["to_name"] = metadata["to_name"].(string)
		sanitizedMetadata["reply_to_email"] = metadata["reply_to_email"].(string)
		sanitizedMetadata["reply_to_name"] = metadata["reply_to_name"].(string)

		notification.Subject = metadata["subject"].(string)
		notification.Metadata = sanitizedMetadata
		// notification.ErrorMessage = nil
	} else {
		return nil, errors.New("metadata is required in notification data")
	}

	err = n.notificationRepository.Create(ctx, notification)
	return notification, err
}

func (n *EmailNotifier) CreateBulkNotifications(ctx context.Context, payload []map[string]any) ([]*models.Notification, error) {
	// Initialize the notification model instance for given payload
	// notifications := []*models.Notification{}
	// for _, p := range payload {
	// 	notification := &models.Notification{}
	// 	notifications = append(notifications, notification)
	// }

	// err := n.notificationRepository.Create(ctx, notifications)
	// return notifications, err
	return nil, nil
}

func (n *EmailNotifier) ChannelType() string {
	return "email"
}