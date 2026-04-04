package emailnotification

import (
	// "log"
	"time"
	"errors"
	"strings"
	"context"

	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
	"github.com/Nitish0007/go_notifier/internal/shared/validators"
	// "github.com/Nitish0007/go_notifier/internal/shared/dto"
	// rabbitmq_utils "github.com/Nitish0007/go_notifier/utils/rabbitmq"
)

type EmailNotificationService struct {
	notificationRepo *EmailNotificationRepository
}

func NewEmailNotificationService(r *EmailNotificationRepository) *EmailNotificationService {
	return &EmailNotificationService{
		notificationRepo: r,
	}
}

func (s *EmailNotificationService) GetNotifications(ctx context.Context, accID int) ([]*EmailNotification, error) {
	list, err := s.notificationRepo.Index(ctx, accID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *EmailNotificationService) CreateEmailCampaign(ctx context.Context, payload *CreateEmailCampaignRequest) (*EmailNotification, error) {
	notificationType, err := StringToEmailNotificationType(payload.Notification.NotificationType)
	if err != nil {
		return nil, err
	}
	notificationStatus, err := StringToEmailNotificationStatus(payload.Notification.Status)
	if err != nil {
		return nil, err
	}

	send_at := sharedhelper.ParseTime(*payload.Notification.SendAt, time.RFC3339)
	if notificationStatus == Scheduled {
		if send_at == nil {
			return nil, errors.New("send_at time is required for scheduled notifications")
		}
		if send_at.Before(*sharedhelper.GetCurrentTime()) {
			return nil, errors.New("send_at time must be in the future")
		}
	}else if strings.ToLower(payload.Notification.Status) == "send_now" {
		if send_at != nil {
			return nil, errors.New("send_at time is not allowed for send_now notifications")
		}
		send_at = sharedhelper.GetCurrentTime()
	}else{
		send_at = nil
	}

	notification := NewEmailNotification(
		payload.Notification.AccountID, 
		payload.Notification.Subject, 
		payload.Notification.Title, 
		notificationType, 
		payload.Notification.ContentID,
		notificationStatus,
		send_at,
	)

	validator := validators.NewModelValidator[EmailNotification]()
	if err := validator.ValidateStruct(notification); err != nil {
		return nil, err
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *EmailNotificationService) CreateEmailTransactionalCampaign(ctx context.Context, payload *CreateEmailTransactionalRequest) (*EmailNotification, error) {
	notificationType, err := StringToEmailNotificationType(payload.Notification.NotificationType)
	if err != nil {
		return nil, err
	}

	notificationStatus, err := StringToEmailNotificationStatus(payload.Notification.Status)
	if err != nil {
		return nil, err
	}

	trans_notif := NewEmailNotification(
		payload.Notification.AccountID, 
		payload.Notification.Subject, 
		payload.Notification.Title, 
		notificationType, 
		payload.Notification.ContentID,
		notificationStatus,
		nil, // send_at should be nil for transactional notifications
	)

	validator := validators.NewModelValidator[EmailNotification]()
	if err := validator.ValidateStruct(trans_notif); err != nil {
		return nil, err
	}

	if err := s.notificationRepo.Create(ctx, trans_notif); err != nil {
		return nil, err
	}

 	return trans_notif, nil
}

func (s *EmailNotificationService) GetNotificationById(ctx context.Context, nID string, accID int) (*EmailNotification, error) {
	n, err := s.notificationRepo.GetByID(ctx, nID, accID)
	if err != nil {
		return nil, err
	}
	return n, nil
}

// func (s *EmailNotificationService) SendOrScheduleNotification(ctx context.Context, n *EmailNotification) error {
// 	// if notification is meant to sent straight away or if its scheduled for next 10 mins then simply push it to queue
// 	if n.Status == Pending && (n.SendAt.Before(time.Now()) || n.SendAt.After(time.Now().Add(10*time.Minute))) {
// 		// push in queue
// 		body := map[string]any{
// 			"notificationID": n.ID,
// 			"accountID":      utils.GetCurrentAccountID(ctx),
// 		}
// 		err := rabbitmq_utils.PushToQueueByName("emailer", rabbitmq_utils.NewJobMessage(body))
// 		if err != nil {
// 			log.Printf("ERROR!!!: %v", err)
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (s *EmailNotificationService) SendNotification(ctx context.Context, notificationID string, accountID int, smtpConfig *dto.SMTPConfiguration) error {
// 	notification, err := s.notificationRepo.GetByID(ctx, notificationID, accountID)
// 	if err != nil {
// 		log.Printf("Error in getting notification: %v", err)
// 		return err
// 	}
// 	log.Printf("Fetched notification by ID => Notification: %v", notification)

// 	if notification == nil {
// 		return errors.New("notification not found")
// 	}

// 	fieldsToUpdate := map[string]any{
// 		"status":        Sent,
// 		"sent_at":       time.Now(),
// 		"error_message": nil,
// 	}
// 	_, err = s.notificationRepo.UpdateNotification(ctx, fieldsToUpdate, notification)
// 	if err != nil {
// 		log.Printf("Error in updating notification: %v", err)
// 		return err
// 	}
// 	return nil
// }
