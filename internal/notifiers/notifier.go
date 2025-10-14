package notifiers

import (
	"context"

	"github.com/Nitish0007/go_notifier/internal/models"
)

type Notifier interface {
	Send(body map[string]any) error
	ChannelType() string
	CreateNotification(ctx context.Context, payload map[string]any) (*models.Notification, error)
	CreateBulkNotifications(ctx context.Context, payload []map[string]any) ([]*models.Notification, error)
}