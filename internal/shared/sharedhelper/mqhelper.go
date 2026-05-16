package sharedhelper

import (
	"log"
	"time"
	"errors"
	"context"
	"encoding/json"
	"github.com/google/uuid"

	"github.com/Nitish0007/go_notifier/internal/common/mq"
)

// JobMetadata is embedded in MQMessage JSON; extra fields help debug retry/DLQ in queue inspectors.
type JobMetadata struct {
	RetryCount int           `json:"retry_count"`
	MaxRetries int           `json:"max_retries"`
	RetryDelay time.Duration `json:"retry_delay"`

	// Failure / routing (set by mq consumer wrapper on handler or decode errors)
	LastError        string    `json:"last_error,omitempty"`
	LastErrorAt      time.Time `json:"last_error_at,omitempty"`
	FailedOnQueue    string    `json:"failed_on_queue,omitempty"`    // queue the consumer was reading
	LastFailureStage string    `json:"last_failure_stage,omitempty"` // "decode" | "handler"
}

type MQMessage struct {
	MessageID  string          `json:"message_id"`
	Payload    json.RawMessage `json:"payload"`
	Metadata   *JobMetadata    `json:"metadata"`
}

func NewMQMessage(payload any, metadata *JobMetadata) (*MQMessage, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal payload: %v", err)
		return nil, err
	}
	return &MQMessage{
		MessageID: uuid.New().String(),
		Payload:   json.RawMessage(payloadBytes),
		Metadata:  metadata,
	}, nil
}

func Encode(m *MQMessage) ([]byte, error) {
	return json.Marshal(m)
}

func Decode(encodedMessage []byte) (*MQMessage, error) {
	m := &MQMessage{}
	err := json.Unmarshal(encodedMessage, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MQMessage) GetMessageID() string {
	return m.MessageID
}

func (m *MQMessage) GetPayload() json.RawMessage {
	return m.Payload
}

func (m *MQMessage) GetMetadata() *JobMetadata {
	return m.Metadata
}

func (m *MQMessage) SetMetadata(metadata *JobMetadata) {
	m.Metadata = metadata
}

func PublishToMQ(ctx context.Context, mqClient mq.MQClient, queueName string, message *MQMessage) error {
	if message == nil {
		return errors.New("message is nil")
	}
	encodedMessage, err := Encode(message)
	if err != nil {
		return err
	}
	return mqClient.Publish(ctx, queueName, encodedMessage)
}

// ConsumeFromMQ forwards to mq.MQClient.Consume. Pass policy nil for legacy behavior (handler error → broker requeue).
// Pass a non-nil mq.ConsumePolicy to enable JSON MQMessage decode, retry/DLQ routing, and optional backoff before retry publish.
func ConsumeFromMQ(ctx context.Context, mqClient mq.MQClient, queueName string, policy *mq.ConsumePolicy, handler mq.MessageHandler) error {
	return mqClient.Consume(ctx, queueName, policy, handler)
}
