package notifier

// import "time"

type NotificationView interface {
	GetRecipient() 			string
	GetSubject() 				string
	GetBody() 					string
	GetHtmlBody() 			string
	GetMetadata() 			map[string]any
	// GetChannel() 				string
	// GetStatus() 				string
	// GetAccountID() 			int
	// GetSendAt() 				time.Time
	// GetSentAt() 				time.Time
	// GetCreatedAt() 			time.Time
	// GetID() 						string
	// GetBatchID() 				string
	// GetJobID() 					string
	// GetErrorMessage() 	string
}