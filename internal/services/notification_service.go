package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Nitish0007/go_notifier/internal/models"
	"github.com/Nitish0007/go_notifier/internal/notifiers"
	"github.com/Nitish0007/go_notifier/internal/repositories"
	"github.com/Nitish0007/go_notifier/utils"
)

type NotificationService struct {
	notifiers map[string]notifiers.Notifier
	notificationRepo *repositories.NotificationRepository
}

func NewNotificationService(list []notifiers.Notifier, nr *repositories.NotificationRepository) *NotificationService {
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
		notificationRepo: nr,
	}
}

func (s *NotificationService) GetNotificationsService(ctx context.Context, accID int) ([]*models.Notification, error) {
	list, err := s.notificationRepo.Index(ctx, accID)
	if err != nil {
		return nil, err
	}
	return list, nil
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

func (s *NotificationService) GetNotificationService(ctx context.Context, nID string, accID int) (*models.Notification, error) {
	n, err := s.notificationRepo.GetByID(ctx, nID, accID)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (s *NotificationService) SendOrScheduleNotification(ctx context.Context, n *models.Notification) error {
	// if notification is meant to sent straight away or if its scheduled for next 10 mins then simply push it to queue
	if n.Status == models.Pending && (n.SendAt.Before(time.Now()) || n.SendAt.After(time.Now().Add(10*time.Minute))) {
		// push in queue
		err := utils.PushToQueue("emailer", n.ID)
		if err != nil {
			log.Printf("ERROR!!!: %v", err)
			return err
		}
	}
	return nil
}