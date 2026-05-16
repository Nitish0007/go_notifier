package workers

import "time"

type Worker interface {
	Run()
	RetryCount() int
	MaxRetries() int
	QueueName() string
	RetryDelay() time.Duration
}
