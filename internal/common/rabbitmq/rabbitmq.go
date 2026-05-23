package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rabbitmq/amqp091-go"

	"github.com/Nitish0007/go_notifier/internal/common/mq"
	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
)

type RabbitMQClient struct {
	conn     *amqp091.Connection
	channel  *amqp091.Channel // used for publishing only
	mu       sync.Mutex
}

var consumerTagSeq atomic.Uint64

func (r *RabbitMQClient) Publish(ctx context.Context, queueName string, message []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	queue, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return errors.New("Failed to declare queue: " + err.Error())
	}

	return r.channel.PublishWithContext(ctx, "", queue.Name, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        message,
		DeliveryMode: amqp091.Persistent,
	})
}

func (r *RabbitMQClient) Consume(ctx context.Context, queueName string, policy *mq.ConsumePolicy, handler mq.MessageHandler) error {
	if handler == nil {
		return errors.New("RabbitMQ handler is required to consume messages")
	}
	if policy != nil && policy.DLQQueue == "" {
		return errors.New("rabbitmq Consume: ConsumePolicy requires DLQQueue")
	}
	if policy != nil && policy.MaxRetries > 0 && policy.RetryQueue == "" {
		return errors.New("rabbitmq Consume: ConsumePolicy requires RetryQueue when MaxRetries > 0")
	}

	wrapped := handler
	if policy != nil {
		wrapped = r.wrapConsumePolicy(ctx, queueName, policy, handler)
	}

	log.Printf("Consuming messages from queue: %s", queueName)

	ch, err := r.conn.Channel()
	if err != nil {
		return errors.New("Failed to open a channel: " + err.Error())
	}
	defer ch.Close()

	if err := ch.Qos(1, 0, false); err != nil {
		return errors.New("Failed to set QoS: " + err.Error())
	}

	consumerTag := fmt.Sprintf("go_notifier-%d-%s", consumerTagSeq.Add(1), queueName)
	msgs, err := ch.Consume(queueName, consumerTag, false, false, false, false, nil)
	if err != nil {
		return errors.New("Failed to consume messages: " + err.Error())
	}

	log.Printf("Consuming messages from queue: %s", queueName)

	for {
		select {
		case <-ctx.Done():
			if cancelErr := ch.Cancel(consumerTag, false); cancelErr != nil {
				log.Printf("rabbitmq Cancel consumer: %v", cancelErr)
			}
			log.Printf("Context cancelled, exiting consumer: %v", ctx.Err())
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed for queue: %s", queueName)
			}

			if err := wrapped(ctx, msg.Body); err != nil {
				log.Printf("Error in processing message: %v", err)
				if nackErr := msg.Nack(false, true); nackErr != nil {
					log.Printf("Failed to nack/requeue message: %v", nackErr)
				}
				continue
			}
			if ackErr := msg.Ack(false); ackErr != nil {
				log.Printf("Failed to acknowledge message: %v", ackErr)
			}
		}
	}
}

const maxMQErrorStringLen = 4096

func truncateMQError(s string) string {
	if len(s) <= maxMQErrorStringLen {
		return s
	}
	return s[:maxMQErrorStringLen] + "…(truncated)"
}

// wrapConsumePolicy decodes MQMessage, runs handler on original bytes, routes failures to
// retry/DLQ with optional delay, then returns nil so the outer loop Ack's the delivery.
func (r *RabbitMQClient) wrapConsumePolicy(ctx context.Context, consumeQueue string, policy *mq.ConsumePolicy, handler mq.MessageHandler) mq.MessageHandler {
	return func(ctx context.Context, body []byte) error {
		msg, err := sharedhelper.Decode(body)
		if err != nil {
			md := &sharedhelper.JobMetadata{
				LastError:        truncateMQError(err.Error()),
				LastErrorAt:      time.Now().UTC(),
				FailedOnQueue:    consumeQueue,
				LastFailureStage: "decode",
			}
			dlq, encErr := sharedhelper.NewMQMessage(map[string]any{
				"error":  "decode_failed",
				"detail": err.Error(),
				"raw":    truncateMQError(string(body)),
			}, md)
			if encErr == nil {
				if pubErr := r.publishMessage(ctx, policy.DLQQueue, dlq); pubErr != nil {
					log.Printf("rabbitmq DLQ publish (decode): %v", pubErr)
				}
			}
			return nil
		}

		if msg.Metadata == nil {
			msg.Metadata = &sharedhelper.JobMetadata{}
		}
		msg.Metadata.MaxRetries = policy.MaxRetries

		handlerErr := handler(ctx, body)
		if handlerErr == nil {
			return nil
		}

		hErr := truncateMQError(handlerErr.Error())
		msg.Metadata.LastError = hErr
		msg.Metadata.LastErrorAt = time.Now().UTC()
		msg.Metadata.FailedOnQueue = consumeQueue
		msg.Metadata.LastFailureStage = "handler"

		msg.Metadata.RetryCount++
		if msg.Metadata.RetryCount > policy.MaxRetries {
			if pubErr := r.publishMessage(ctx, policy.DLQQueue, msg); pubErr != nil {
				log.Printf("rabbitmq DLQ publish: message_id=%s queue=%s err=%v", msg.MessageID, consumeQueue, pubErr)
			} else {
				log.Printf("rabbitmq DLQ: message_id=%s from=%s retries=%d max=%d last_error=%q",
					msg.MessageID, consumeQueue, msg.Metadata.RetryCount, policy.MaxRetries, hErr)
			}
			return nil
		}

		if policy.MaxRetries > 0 && policy.RetryQueue != "" {
			if policy.BaseDelay > 0 {
				delay := policy.BaseDelay * time.Duration(msg.Metadata.RetryCount)
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(delay):
				}
			}
			if pubErr := r.publishMessage(ctx, policy.RetryQueue, msg); pubErr != nil {
				log.Printf("rabbitmq retry publish: message_id=%s queue=%s err=%v", msg.MessageID, consumeQueue, pubErr)
			} else {
				log.Printf("rabbitmq retry: message_id=%s from=%s attempt=%d/%d last_error=%q next=%s",
					msg.MessageID, consumeQueue, msg.Metadata.RetryCount, policy.MaxRetries, hErr, policy.RetryQueue)
			}
		}
		return nil
	}
}

func (r *RabbitMQClient) publishMessage(ctx context.Context, queue string, msg *sharedhelper.MQMessage) error {
	body, err := sharedhelper.Encode(msg)
	if err != nil {
		return err
	}
	return r.publishBytes(ctx, queue, body)
}

func (r *RabbitMQClient) publishBytes(ctx context.Context, queue string, body []byte) error {
	return r.Publish(ctx, queue, body)
}


const (
	dialMaxAttempts = 15
	dialInterval    = 2 * time.Second
)
var rabbitMQ *RabbitMQClient
var clientMutex  sync.Mutex

func GetRabbitMQURI() string {
	return os.Getenv("RABBITMQ_URL")
}


func dialRabbitMQ(uri string) (*amqp091.Connection, error) {
	var last error

	for attempt := 1; attempt <= dialMaxAttempts; attempt++ {
		con, err := amqp091.Dial(uri)
		if err == nil {
			return con, nil
		}

		last = err

		log.Printf("RabbitMQ dial attempt %d/%d failed: %v",
			attempt, dialMaxAttempts, err)

		if attempt < dialMaxAttempts {
			time.Sleep(dialInterval)
		}
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w",
		dialMaxAttempts, last)
}

func NewRabbitMQClient() (mq.MQClient, error) {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	if rabbitMQ != nil {
		return rabbitMQ, nil
	}

	uri := GetRabbitMQURI()
	if uri == "" {
		return nil, errors.New("RABBITMQ_URL is not set")
	}
	
	con, err := dialRabbitMQ(uri)
	if err != nil {
		return nil, err
	}

	ch, err := con.Channel()
	if err != nil {
		_ = con.Close()
		return nil, errors.New("Failed to open a channel: " + err.Error())
	}

	rabbitMQ = &RabbitMQClient{conn: con, channel: ch}
	log.Println("Connected to RabbitMQ successfully")
	return rabbitMQ, nil
}



