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

	if metadata, exists := data["metadata"].(map[string]any); exists {
		if len(metadata) == 0 {
			return nil, errors.New("metadata can't be empty")
		}
		if _, ok := metadata["from_email"]; !ok || metadata["from_email"] == "" {
			return nil, errors.New("from_email is required in metadata")
		}
		if _, ok := metadata["from_name"]; !ok || metadata["from_name"] == "" {
			return nil, errors.New("from_name is required in metadata")
		}
		if _, ok := metadata["to_email"]; !ok || metadata["to_email"] == "" {
			return nil, errors.New("to_email is required in metadata")
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
	} else {
		return nil, errors.New("metadata is required in notification data")
	}

	// adding to_email in metadata
	metadata := data["metadata"].(map[string]any)
	metadata["to_email"] = data["to"]
	// removing from data
	delete(data, "to")

	if data["message"] == nil || data["message"] == "" {
		return nil, errors.New("message is required in notification data")
	}
	if data["body"] == nil || data["body"] == "" {
		return nil, errors.New("body is required in notification data")
	}
	if data["content"] == nil || data["content"] == "" {
		return nil, errors.New("content is required in notification data")
	}
	if data["html_part"] == nil || data["html_part"] == "" {
		return nil, errors.New("html_part is required in notification data")
	}
	// Validation code ends here


	notifierObj := s.notifiers[notificationType]
	log.Printf(">>>>>>>>>>>>> Notifier Object: %+v", notifierObj)
	n, err := notifierObj.CreateNotification(ctx, data)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, errors.New("Failed to create notification")
	}
	
	return n, nil
}

func (s *NotificationService) CreateBulkNotifications(ctx context.Context, data []map[string]any) ([]*models.Notification, error) {

	return nil, nil
}

func (s *NotificationService) SendOrScheduleNotification(ctx context.Context, n *models.Notification) error {

	return nil
}