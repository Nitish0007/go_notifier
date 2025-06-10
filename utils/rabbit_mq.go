package utils

import (
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

// func CreateQueue(conn *rbmq.Connection, queue_name string) error {

// }