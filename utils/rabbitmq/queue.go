package rabbitmq_utils

import (
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

func (q *Queue) PushToMain(body map[string]any) error {
	err := PushToQueue(q.Main.Name, body)
	if err != nil {
		return err
	}
	return nil
}

func (q *Queue) PushToRetry(body map[string]any) error {
	err := PushToQueue(q.Retry.Name, body)
	if err != nil {
		return err
	}
	return nil
}

func (q *Queue) PushToDLQ(body map[string]any) error {
	err := PushToQueue(q.DLQ.Name, body)
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
	retryQueue, err := CreateQueue(ch, queue_name + "_retry")
	if err != nil {
		return nil, err
	}

	// create dead letter queue
	dlq, err := CreateQueue(ch, queue_name + "_dlq")
	if err != nil {
		return nil, err
	}

	return &Queue{
		Main:  mainQueue,
		Retry: retryQueue,
		DLQ:   dlq,
	}, nil
}