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
	rabbitmq_utils "github.com/Nitish0007/go_notifier/utils/rabbitmq"
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

func (s *NotificationService) GetNotifications(ctx context.Context, accID int) ([]*models.Notification, error) {
	list, err := s.notificationRepo.Index(ctx, accID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *NotificationService) CreateNotification(ctx context.Context, data map[string]any) (*models.Notification, error) {
	// Validate the notification data
	validPayload, err := utils.ValidateNotificationPayload(data)
	if err != nil {
		return nil, err
	}

	notificationType, _ := validPayload["channel"].(string)
	notifierObj := s.notifiers[notificationType]
	n, err := notifierObj.CreateNotification(ctx, validPayload)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, errors.New("failed to create notification")
	}
	
return n, nil
}

func (s *NotificationService) GetNotificationById(ctx context.Context, nID string, accID int) (*models.Notification, error) {
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
		body := map[string]any{
			"notificationID": n.ID,
			"accountID": utils.GetCurrentAccountID(ctx),
		}
		err := rabbitmq_utils.PushToQueue("emailer", body)
		if err != nil {
			log.Printf("ERROR!!!: %v", err)
			return err
		}
	}
	return nil
}