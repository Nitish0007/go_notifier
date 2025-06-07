package handlers

import (
	"github.com/Nitish0007/go_notifier/internal/services"
)

type NotificationHandler struct{
	notificationService *services.NotificationService
}

func NewNotificationHandler(s *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: s,
	}
}