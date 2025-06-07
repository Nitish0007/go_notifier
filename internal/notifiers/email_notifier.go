package notifiers

import (
	"github.com/Nitish0007/go_notifier/internal/repositories"
)

type EmailNotifier struct {
	notificationRepository *repositories.NotificationRepository
};

func NewEmailNotifier(r *repositories.NotificationRepository) *EmailNotifier {
	return &EmailNotifier{notificationRepository: r}
}

func (n *EmailNotifier) Send(to string, payload map[string]any) error {
	
	return nil
}

func (n *EmailNotifier) ChannelType() string {
	return "email"
}