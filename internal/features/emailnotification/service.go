package emailnotification

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
	"github.com/Nitish0007/go_notifier/internal/shared/validators"
)

type EmailNotificationService struct {
	notificationRepo *EmailNotificationRepository
}

func NewEmailNotificationService(r *EmailNotificationRepository) *EmailNotificationService {
	return &EmailNotificationService{
		notificationRepo: r,
	}
}

func (s *EmailNotificationService) GetNotifications(ctx context.Context, accID int64) ([]*EmailNotification, error) {
	list, err := s.notificationRepo.Index(ctx, accID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *EmailNotificationService) CreateEmailCampaign(ctx context.Context, payload *CreateEmailCampaignRequest) (*EmailCampaignResponse, error) {
	typeStr := strings.TrimSpace(payload.Notification.NotificationType)
	if typeStr == "" {
		typeStr = "campaign"
	}
	notificationType, err := StringToEmailNotificationType(typeStr)
	if err != nil {
		return nil, err
	}
	if notificationType != Campaign {
		return nil, errors.New("invalid notification type: must be campaign")
	}

	statusStr := strings.ToLower(strings.TrimSpace(payload.Notification.Status))
	var sendAt *time.Time
	var dbStatus EmailNotificationStatus

	switch statusStr {
	case "scheduled":
		if payload.Notification.SendAt == nil {
			return nil, errors.New("send_at time is required for scheduling notifications")
		}
		t := sharedhelper.ParseTime(*payload.Notification.SendAt, time.RFC3339)
		if t == nil {
			return nil, errors.New("invalid send_at time format, should be in RFC3339 format")
		}
		if !t.After(*sharedhelper.GetCurrentTime()) {
			return nil, errors.New("send_at time must be in the future")
		}
		sendAt = t
		dbStatus = Scheduled
	case "send_now":
		sendAt = sharedhelper.GetCurrentTime()
		dbStatus = Scheduled
	case "draft":
		sendAt = nil
		dbStatus = Draft
	default:
		return nil, errors.New("invalid status: must be draft, scheduled, or send_now")
	}

	if len(payload.Notification.ListIDs) == 0 {
		return nil, errors.New("list_ids are required")
	}

	notification := NewEmailNotification(
		payload.Notification.AccountID,
		payload.Notification.Subject,
		payload.Notification.Title,
		notificationType,
		payload.Notification.ContentID,
		dbStatus,
		sendAt,
		payload.Notification.FromName,
		payload.Notification.FromEmail,
		payload.Notification.ReplyToEmail,
	)

	validator := validators.NewModelValidator[EmailNotification]()
	if err := validator.ValidateStruct(notification); err != nil {
		return nil, err
	}

	if err := s.notificationRepo.CreateCampaignWithList(ctx, notification, payload.Notification.ListIDs); err != nil {
		return nil, err
	}

	return &EmailCampaignResponse{
		ID: notification.ID,
		AccountID: notification.AccountID,
		Subject: notification.Subject,
		Title: notification.Title,
		FromName: notification.FromName,
		FromEmail: notification.FromEmail,
		ReplyToEmail: notification.ReplyToEmail,
		ContentID: notification.ContentID,
		ListIDs: payload.Notification.ListIDs,
		NotificationType: string(notification.NotificationType),
		Status: string(notification.Status),
		SendAt: notification.SendAt,
		SentAt: notification.SentAt,
		CreatedAt: notification.CreatedAt,
		// UpdatedAt: notification.UpdatedAt,
	}, nil
}

func (s *EmailNotificationService) CreateEmailTransactionalCampaign(ctx context.Context, payload *CreateEmailTransactionalRequest) (*EmailNotification, error) {
	typeStr := strings.TrimSpace(payload.Notification.NotificationType)
	if typeStr == "" {
		typeStr = "transactional"
	}
	notificationType, err := StringToEmailNotificationType(typeStr)
	if err != nil {
		return nil, err
	}
	if notificationType != Transactional {
		return nil, errors.New("transactional creation requires notification_type transactional")
	}

	statusStr := strings.TrimSpace(payload.Notification.Status)
	if statusStr == "" {
		statusStr = "trans"
	}
	notificationStatus, err := StringToEmailNotificationStatus(statusStr)
	if err != nil {
		return nil, err
	}

	transNotif := NewEmailNotification(
		payload.Notification.AccountID,
		payload.Notification.Subject,
		payload.Notification.Title,
		notificationType,
		payload.Notification.ContentID,
		notificationStatus,
		nil, // send_at is nil for transactional notifications
		payload.Notification.FromName,
		payload.Notification.FromEmail,
		payload.Notification.ReplyToEmail,
	)

	validator := validators.NewModelValidator[EmailNotification]()
	if err := validator.ValidateStruct(transNotif); err != nil {
		return nil, err
	}

	if err := s.notificationRepo.Create(ctx, transNotif); err != nil {
		return nil, err
	}

	return transNotif, nil
}

func (s *EmailNotificationService) GetNotificationById(ctx context.Context, nID int64, accID int64) (*EmailNotification, error) {
	n, err := s.notificationRepo.GetByID(ctx, nID, accID)
	if err != nil {
		return nil, err
	}
	return n, nil
}

// func (s *EmailNotificationService) SendOrScheduleNotification(ctx context.Context, n *EmailNotification) error {
// 	// if notification is meant to sent straight away or if its scheduled for next 10 mins then simply push it to queue
// 	if n.Status == Trans && (n.SendAt.Before(time.Now()) || n.SendAt.After(time.Now().Add(10*time.Minute))) {
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
