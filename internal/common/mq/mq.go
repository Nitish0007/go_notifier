package mq

import "context"

type MessageHandler func(ctx context.Context, message []byte) error

type MQClient interface {
	Publish(ctx context.Context, queueName string, message []byte) error
	// Consume subscribes to queueName. If policy is nil, handler receives raw message bytes;
	// on handler error the message is Nack+requeued (legacy). If policy is non-nil, messages
	// are treated as sharedhelper.MQMessage JSON: handler still receives the same body bytes;
	// failures are routed to RetryQueue/DLQ with optional delay, then the delivery is Ack'd.
	Consume(ctx context.Context, queueName string, policy *ConsumePolicy, handler MessageHandler) error
}