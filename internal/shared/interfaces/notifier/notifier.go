package notifier

import "context"

// Notifier sends a message on one channel (Strategy).
type Notifier interface {
	Channel() Channel
	Notify(ctx context.Context, req DeliveryRequest, cfg ProviderConfig) error
}
