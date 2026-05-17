package notifier

// Channel identifies a delivery channel (email, sms, etc.).
type Channel string

const (
	ChannelEmail Channel = "email"
)
