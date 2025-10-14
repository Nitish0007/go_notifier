package utils

import (
	"encoding/json"
	"log"

	rbmq "github.com/rabbitmq/amqp091-go"
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
	return ch, err
}

func CreateQueue(ch *rbmq.Channel, queue_name string) (*rbmq.Queue, error) {
	q, err := ch.QueueDeclare(
		queue_name, // name
		false,   		// durable
		false,   		// delete when unused
		false,   		// exclusive
		false,   		// no-wait
		nil,     		// arguments
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
			ContentType: "application/json",
			Body: jsonBody,			
		},
	)

	if err != nil {
		failOnError(err, "Failed while publishing")
		return err
	}

	return nil
}