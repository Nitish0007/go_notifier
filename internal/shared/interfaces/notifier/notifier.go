package notifier

import (
	"context"

	"github.com/Nitish0007/go_notifier/internal/shared/dto"
)

type Notifier interface {
	Notify(notification NotificationView, smtpConfig *dto.SMTPConfiguration) error
	ChannelType() string
	CreateNotification(ctx context.Context, payload map[string]any) (NotificationView, error)
	// CreateBulkNotifications(ctx context.Context, payload []map[string]any) ([]NotificationView, error)
}