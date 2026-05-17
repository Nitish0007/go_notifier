package notifier

// ProviderConfig is channel-specific credentials/settings (SMTP, Twilio, etc.).
type ProviderConfig interface {
	Channel() Channel
}
