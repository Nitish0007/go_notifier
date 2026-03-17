package notification

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Nitish0007/go_notifier/utils"
	"github.com/Nitish0007/go_notifier/internal/shared/dto"
	notifierInterface "github.com/Nitish0007/go_notifier/internal/shared/interfaces/notifier"
	rabbitmq_utils "github.com/Nitish0007/go_notifier/utils/rabbitmq"
)

type NotificationService struct {
	notifiers        map[string]notifierInterface.Notifier
	notificationRepo *NotificationRepository
}

func NewNotificationService(list []notifierInterface.Notifier, nr *NotificationRepository) *NotificationService {
	nList := make(map[string]notifierInterface.Notifier)
	for _, val := range list {
		if val.ChannelType() == "email" || val.ChannelType() == "sms" || val.ChannelType() == "in_app" {
			nList[val.ChannelType()] = val
		} else {
			log.Printf(">>>>>>>>>>>>> Unknown Notifier: %v", val.ChannelType())
		}
	}
	return &NotificationService{
		notifiers:        nList,
		notificationRepo: nr,
	}
}

func (s *NotificationService) GetNotifications(ctx context.Context, accID int) ([]*Notification, error) {
	list, err := s.notificationRepo.Index(ctx, accID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *NotificationService) CreateNotification(ctx context.Context, data map[string]any) (*Notification, error) {
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
	notif, ok := n.(*Notification)
	if !ok {
		return nil, errors.New("failed to cast to Notification")
	}

	return notif, nil
}

func (s *NotificationService) GetNotificationById(ctx context.Context, nID string, accID int) (*Notification, error) {
	n, err := s.notificationRepo.GetByID(ctx, nID, accID)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (s *NotificationService) SendOrScheduleNotification(ctx context.Context, n *Notification) error {
	// if notification is meant to sent straight away or if its scheduled for next 10 mins then simply push it to queue
	if n.Status == Pending && (n.SendAt.Before(time.Now()) || n.SendAt.After(time.Now().Add(10*time.Minute))) {
		// push in queue
		body := map[string]any{
			"notificationID": n.ID,
			"accountID":      utils.GetCurrentAccountID(ctx),
		}
		err := rabbitmq_utils.PushToQueueByName("emailer", rabbitmq_utils.NewJobMessage(body))
		if err != nil {
			log.Printf("ERROR!!!: %v", err)
			return err
		}
	}
	return nil
}

func (s *NotificationService) SendNotification(ctx context.Context, notificationID string, accountID int, smtpConfig *dto.SMTPConfiguration) error {
	notification, err := s.notificationRepo.GetByID(ctx, notificationID, accountID)
	if err != nil {
		log.Printf("Error in getting notification: %v", err)
		return err
	}
	log.Printf("Fetched notification by ID => Notification: %v", notification)

	if notification == nil {
		return errors.New("notification not found")
	}

	channelString, err := ChannelToString(notification.Channel)
	if err != nil {
		log.Printf("Error in getting channel string: %v", err)
		return err
	}
	notifier := s.notifiers[channelString]
	if notifier == nil {
		return errors.New("no notifier allocated for channel type: " + channelString)
	}

	err = notifier.Notify(notification, smtpConfig)
	if err != nil {
		log.Printf("Error in sending notification: %v", err)
		fieldsToUpdate := map[string]any{
			"status":        Failed,
			"sent_at":       time.Now(),
			"error_message": err.Error(),
		}
		_, err = s.notificationRepo.UpdateNotification(ctx, fieldsToUpdate, notification)
		if err != nil {
			log.Printf("Error in updating notification: %v", err)
			return err
		}
		return err
	}

	fieldsToUpdate := map[string]any{
		"status":        Sent,
		"sent_at":       time.Now(),
		"error_message": nil,
	}
	_, err = s.notificationRepo.UpdateNotification(ctx, fieldsToUpdate, notification)
	if err != nil {
		log.Printf("Error in updating notification: %v", err)
		return err
	}
	return nil
}
