package services

import (
	"context"
	"errors"
	"log"

	"github.com/Nitish0007/go_notifier/internal/models"
	"github.com/Nitish0007/go_notifier/internal/notifiers"
	"github.com/Nitish0007/go_notifier/utils"
)

type NotificationService struct {
	notifiers map[string]notifiers.Notifier
}

func NewNotificationService(list []notifiers.Notifier) *NotificationService {
	nList := make(map[string]notifiers.Notifier)
	for _, val := range list {
		if val.ChannelType() == "email" || val.ChannelType() == "sms" || val.ChannelType() == "in_app" {
			nList[val.ChannelType()] = val
		}else{
			log.Printf(">>>>>>>>>>>>> Unknown Notifier: %v", val.ChannelType())
		}
	}
	return &NotificationService{
		notifiers: nList,
	}
}

func (s *NotificationService) CreateNotification(ctx context.Context, data map[string]any) (*models.Notification, error) {
	// Validate the notification data
	notificationType, exists := data["channel"].(string)
	if !exists {
		return nil, errors.New("channel is required in notification data")
	}

	if !utils.IsValidChannelType(notificationType) {
		return nil, errors.New("invalid channel type provided")
	}

	if recipient, exists := data["to"].(string); !exists || recipient == "" {
		return nil, errors.New("recipient 'to' is required in notification data")
	}

	body, exists := data["body"]
	if !exists || body == nil || body == "" {
		data["body"] = ""
	}
	htmlBody, exists := data["html_body"]
	if !exists || htmlBody == nil || htmlBody == "" {
		data["html_body"] = ""
	}

	if data["body"] == "" && data["html_body"] == "" {
		return nil, errors.New("either 'body' or 'html_body' must be provided in notification data")
	}
	// Validation code ends here

	notifierObj := s.notifiers[notificationType]
	n, err := notifierObj.CreateNotification(ctx, data)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, errors.New("failed to create notification")
	}
	
	return n, nil
}

func (s *NotificationService) CreateBulkNotifications(ctx context.Context, data []map[string]any) ([]*models.Notification, error) {

	return nil, nil
}

func (s *NotificationService) SendOrScheduleNotification(ctx context.Context, n *models.Notification) error {

	return nil
}