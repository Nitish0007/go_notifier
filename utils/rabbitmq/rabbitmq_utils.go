package rabbitmq_utils

import (
	"log"
	"time"

	rbmq "github.com/rabbitmq/amqp091-go"
)

const (
	MAX_RETRIES = 5
	RETRY_DELAY = 1 * time.Minute
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func ConnectMQ() *rbmq.Connection {
	conn, err := rbmq.Dial("amqp://user:password@notifier_rbmq:5672/")
	failOnError(err, "Failed to connect to Rabbit MQ!")
	return conn
}

func CreateChannel(conn *rbmq.Connection) (*rbmq.Channel, error) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a Channel")

	// set QoS (Quality of Service) for the channel
	err = ch.Qos(1, 0, false)
	if err != nil {
		failOnError(err, "Failed to set QoS")
		return nil, err
	}
	return ch, nil
}

func CreateChannelWithQos(conn *rbmq.Connection, prefetchCount int, prefetchSize int, global bool) (*rbmq.Channel, error) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a Channel")

	err = ch.Qos(prefetchCount, prefetchSize, global)
	if err != nil {
		failOnError(err, "Failed to set QoS")
		return nil, err
	}
	return ch, nil
}

func CreateQueue(ch *rbmq.Channel, queue_name string) (*rbmq.Queue, error) {
	q, err := ch.QueueDeclare(
		queue_name, // name
		true,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to Declare queue")
	return &q, err
}

func PushToQueue(queue *rbmq.Queue, jobMessage *JobMessage) error {
	jsonBody, err := jobMessage.ToJSON()
	failOnError(err, "Error converting job message to JSON")

	conn := ConnectMQ()
	defer conn.Close()

	ch, _ := CreateChannel(conn)
	defer ch.Close()

	err = ch.Publish(
		"",
		queue.Name,
		false,
		false,
		rbmq.Publishing{
			DeliveryMode: rbmq.Persistent,
			ContentType: "application/json",
			Body:        jsonBody,
		},
	)

	if err != nil {
		failOnError(err, "Failed while publishing")
		return err
	}

	return nil
}

func PushToQueueByName(queue_name string, jobMessage *JobMessage) error {
	conn := ConnectMQ()
	defer conn.Close()

	ch, _ := CreateChannel(conn)
	defer ch.Close()

	q, err := CreateQueue(ch, queue_name)
	if err != nil {
		return err
	}
	return PushToQueue(q, jobMessage)
}

func ReadFromQueue(q *rbmq.Queue, ch *rbmq.Channel) (<-chan *JobMessage, error) {
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		failOnError(err, "Failed to consume messages")
		return nil, err
	}

	// create a channel to store the job messages
	jobMessages := make(chan *JobMessage, len(msgs))
	for msg := range msgs {
		jobMsg := NewJobMessage(map[string]any{})
		err := jobMsg.FromJSON(msg.Body)
		if err != nil {
			failOnError(err, "Failed to unmarshal job message")
			return nil, err
		}
		jobMessages <- jobMsg
	}
	close(jobMessages)

	return jobMessages, nil
}

