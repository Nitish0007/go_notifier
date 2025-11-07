package workers

import (
	"log"

	"encoding/json"

	"github.com/Nitish0007/go_notifier/internal/notifiers"
	"github.com/Nitish0007/go_notifier/internal/repositories"
	"github.com/Nitish0007/go_notifier/utils"
	rabbitmq_utils "github.com/Nitish0007/go_notifier/utils/rabbitmq"
)

func ConsumeEmailNotifications() {
	conn := rabbitmq_utils.ConnectMQ()
	defer conn.Close()

	ch, _ := rabbitmq_utils.CreateChannel(conn)
	defer ch.Close()

	queue_name := "emailer"
	q, err := rabbitmq_utils.CreateQueue(ch, queue_name)
	if err != nil {
		log.Printf("Error in worker: %v", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Error in worker: %v", err)
		return
	}

	log.Println("[*] Waiting for Job in queue. Press Ctrl+C to exit")

	dbConn, err := utils.ConnectDB()
	if err != nil {
		log.Printf("Error in db connection: %v", err)
		return
	}
	notificationRepo := repositories.NewNotificationRepository(dbConn)
	emailNotifier := notifiers.NewEmailNotifier(notificationRepo)
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var body map[string]any
			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				log.Printf("Unable to decode json: %v", err)
				continue
			}

			err = emailNotifier.Send(body)
			if err != nil {
				log.Printf("Error sending email: %v", err)
			}
		}
	}()

	<-forever
}