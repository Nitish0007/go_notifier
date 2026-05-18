package notifier

import (
	"fmt"
	"context"

	notifierif "github.com/Nitish0007/go_notifier/internal/shared/interfaces/notifier"
)

// Registry holds channel notifiers (Strategy registry).
type Registry struct {
	byChannel map[notifierif.Channel]notifierif.Notifier
}

func NewRegistry() *Registry {
	return &Registry{byChannel: make(map[notifierif.Channel]notifierif.Notifier)}
}

func (r *Registry) Register(n notifierif.Notifier) {
	r.byChannel[n.Channel()] = n
}

func (r *Registry) Notify(ctx context.Context, ch notifierif.Channel, req notifierif.DeliveryRequest, cfg notifierif.ProviderConfig) error {
	n, ok := r.byChannel[ch]
	if !ok {
		return fmt.Errorf("no notifier registered for channel %q", ch)
	}
	if cfg != nil && cfg.Channel() != ch {
		return fmt.Errorf("config channel %q does not match %q", cfg.Channel(), ch)
	}
	return n.Notify(ctx, req, cfg)
}
