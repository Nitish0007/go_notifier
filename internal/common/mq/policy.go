package mq

import "time"

// NoConsumePolicy is nil — use as the policy argument to disable retry/DLQ wrapping.
var NoConsumePolicy *ConsumePolicy

// ConsumePolicy enables decode → handler → retry queue / DLQ / backoff inside Consume.
// Pass nil to Consume for legacy behavior: handler receives raw bytes; on error, Nack+requeue.
//
// MaxRetries: number of failures allowed *after* the first attempt before sending to DLQ.
//   - 0: never use RetryQueue; first handler error publishes to DLQ only.
//   - N: after N failures (RetryCount > N), publish to DLQ; otherwise publish to RetryQueue.
//
// BaseDelay: multiplied by RetryCount before republishing to RetryQueue (0 = no wait).
// Queues: RetryQueue and DLQQueue must be set when policy is non-nil (MainQueue is optional, for docs).
type ConsumePolicy struct {
	MaxRetries   int
	BaseDelay    time.Duration
	MainQueue    string
	RetryQueue   string
	DLQQueue     string
}
