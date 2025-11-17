package notifiers

import (
	"context"

	"github.com/Nitish0007/go_notifier/internal/models"
)

type Notifier interface {
	// here cnoifguration must be generic and should be able to handle any configuration
	Send(notification *models.Notification, smtpConfig *models.SMTPConfiguration) error
	ChannelType() string
	CreateNotification(ctx context.Context, payload map[string]any) (*models.Notification, error)
	CreateBulkNotifications(ctx context.Context, payload []map[string]any) ([]*models.Notification, error)
}
