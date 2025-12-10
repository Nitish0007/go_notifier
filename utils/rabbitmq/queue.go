package rabbitmq_utils

import (
	"context"
	"time"

	rbmq "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	Main  *rbmq.Queue // main queue
	Retry *rbmq.Queue // retry queue
	DLQ   *rbmq.Queue // dead letter queue
}

func NewQueue(queue_name string) (*Queue, error) {
	q, err := setupQueue(queue_name)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (q *Queue) PushToMain(jobMessage *JobMessage) error {
	err := PushToQueue(q.Main, jobMessage)
	if err != nil {
		return err
	}
	return nil
}

func (q *Queue) PushToRetry(jobMessage *JobMessage) error {
	jid := jobMessage.GetJobID()
	metadata, err := GetJobMetadata(context.Background(), jid)
	if err != nil {
		return err
	}

	// create new metadata if not already present
	if metadata == nil {
		metadata = NewJobMetadata(1, MAX_RETRIES, CalculateRetryDelay(1))
		err = StoreJobMetadata(context.Background(), jid, *metadata)
		if err != nil {
			return err
		}
		return PushToQueue(q.Retry, jobMessage)
	}

	// increment retry count
	metadata.RetryCount++

	if metadata.RetryCount > MAX_RETRIES {
		return q.PushToDLQ(jobMessage)
	}
	// calculate retry delay
	metadata.RetryDelay = CalculateRetryDelay(metadata.RetryCount)
	metadata.MaxRetries = MAX_RETRIES
	// store updated metadata
	err = StoreJobMetadata(context.Background(), jid, *metadata)
	if err != nil {
		return err
	}
	err = PushToQueue(q.Retry, jobMessage)
	if err != nil {
		return err
	}
	return nil
}

func (q *Queue) PushToDLQ(jobMessage *JobMessage) error {
	err := PushToQueue(q.DLQ, jobMessage)
	if err != nil {
		return err
	}
	return nil
}

// private function to setup the queue
func setupQueue(queue_name string) (*Queue, error) {
	conn := ConnectMQ()
	defer conn.Close()

	ch, err := CreateChannel(conn)
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	// create main queue
	mainQueue, err := CreateQueue(ch, queue_name)
	if err != nil {
		return nil, err
	}

	// create retry queue
	retryQueue, err := CreateQueue(ch, queue_name+"_retry")
	if err != nil {
		return nil, err
	}

	// create dead letter queue
	dlq, err := CreateQueue(ch, queue_name+"_dlq")
	if err != nil {
		return nil, err
	}

	return &Queue{
		Main:  mainQueue,
		Retry: retryQueue,
		DLQ:   dlq,
	}, nil
}

func CalculateRetryDelay(retryNumber int) time.Duration {
	return time.Duration(retryNumber) * RETRY_DELAY
}


