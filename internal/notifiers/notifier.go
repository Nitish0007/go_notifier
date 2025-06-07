package notifiers

type Notifier interface {
	Send(to string, payload map[string]any) error
	ChannelType() string
}