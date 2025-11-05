package workers

import (
	"encoding/json"
	"log"

	"github.com/Nitish0007/go_notifier/utils"
)

const (
	NOTIFICATION_CREATE_QUEUE = "notification_batch"
)

func ConsumeNotificationBatch() {
	conn := utils.ConnectMQ()
	defer conn.Close()

	ch, err := utils.CreateChannel(conn)
	if err != nil {
		log.Printf("Error creating channel")
		return
	}
	defer ch.Close()

	queue_name := NOTIFICATION_CREATE_QUEUE
	q, err := utils.CreateQueue(ch, queue_name)
	if err != nil {
		log.Printf("Error in worker: %v", err)
		return
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("Error in worker: %v", err)
		return
	}

	log.Printf("messages present in batch_queue: %v\n", msgs)
	log.Println("[*] Waiting for Job in queue. Press Ctrl+C to exit")

	forever := make(chan bool)
	go func() {
		log.Println("Goroutine started, waiting for messages...")
		for d := range msgs {
			log.Printf("======>> message : %v", d)
			var body map[string]any
			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				log.Printf("Error in unmarshalling body: %v", err)
				continue
			}

		}
	}()

	<-forever
}
