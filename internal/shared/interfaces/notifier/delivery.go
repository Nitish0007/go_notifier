package notifier

// DeliveryRequest is channel-neutral input for notifiers.
type DeliveryRequest struct {
	AccountID int64
	Recipient string
	Subject   string
	Body      string
	HTMLBody  string
	Metadata  map[string]string
}
