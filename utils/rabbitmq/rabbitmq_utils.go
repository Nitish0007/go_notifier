package rabbitmq_utils

import (
	"encoding/json"
	"log"
	"time"
	// "context"

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

func PushToQueue(queue_name string, body map[string]any) error {
	jsonBody, err := json.Marshal(body)
	failOnError(err, "Error converting body to JSON")

	conn := ConnectMQ()
	defer conn.Close()

	ch, _ := CreateChannel(conn)
	defer ch.Close()

	q, err := CreateQueue(ch, queue_name)
	failOnError(err, "Error creating Queue")
	err = ch.Publish(
		"",
		q.Name,
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

// func ProcessMsgWithRetry(ch *rbmq.Channel, queue *rbmq.Queue, funcToExecute func(context.Context, map[string]any) error) error {
// 	msgs, err := ch.Consume(
// 		queue.Name,
// 		"",		// consumer tag
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)

// 	if err != nil {
// 		failOnError(err, "Failed to consume messages")
// 		return err
// 	}

// 	ctx := context.Background()

// 	for retryCount := 1; retryCount <= MAX_RETRIES; retryCount++ {
// 		select {
// 		case msg := <-msgs:
// 			var body map[string]any
// 			err = json.Unmarshal(msg.Body, &body)
			
// 			if err != nil {
// 				log.Printf("Error unmarshalling body: %v", err)
// 				msg.Ack(false)
// 				continue
// 			}

// 			err = funcToExecute(ctx, body)
// 			if err == nil {
// 				msg.Ack(false)
// 				return nil
// 			}
// 			log.Printf("Error executing function: %v", err)
// 			time.Sleep(CalculateRetryDelay(retryCount))
// 		default:
// 			log.Printf("No message received after %d retries", MAX_RETRIES)
// 			return nil
// 		}
// 	}

// 	return nil
// }

func CalculateRetryDelay(retryNumber int) time.Duration {
	return time.Duration(retryNumber) * RETRY_DELAY
}
