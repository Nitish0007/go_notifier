package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nitish0007/go_notifier/internal/notifiers"
	"github.com/Nitish0007/go_notifier/internal/repositories"
	"github.com/Nitish0007/go_notifier/internal/services"
	"github.com/Nitish0007/go_notifier/internal/workers"
	"github.com/Nitish0007/go_notifier/utils"
	rabbitmq_utils "github.com/Nitish0007/go_notifier/utils/rabbitmq"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile) // configuring logger to print filename and line number
	log.Println("\n\nStarting Workers...")

	// make database connection for workers
	dbConn, err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// make RabbitMQ connection for workers
	rbmqConn := rabbitmq_utils.ConnectMQ()
	defer rbmqConn.Close()

	// create context for workers
	ctx := context.Background()

	// Initialize Repositories by injecting db connection dependency
	notificationRepo := repositories.NewNotificationRepository(dbConn)
	notificationBatchRepo := repositories.NewNotificationBatchRepository(dbConn)
	notificationBatchErrorRepo := repositories.NewNotificationBatchErrorRepo(dbConn)

	// initialize notifiers
	emailNotifier := notifiers.NewEmailNotifier(notificationRepo)

	// Initialize Services by injecting corresponding repository dependency
	notificationService := services.NewNotificationService([]notifiers.Notifier{emailNotifier}, notificationRepo)
	blkNotificationService := services.NewBulkNotificationService(notificationRepo, notificationBatchRepo, notificationBatchErrorRepo)

	// Initialize workers by injecting dependencies
	notificationBatchWorker := workers.NewNotificationBatchWorker(dbConn, rbmqConn, ctx, blkNotificationService)
	schedulerWorker := workers.NewNotificationSchedulerWorker(dbConn, rbmqConn, ctx, blkNotificationService)
	emailWorker := workers.NewEmailWorker(dbConn, rbmqConn, ctx, notificationService)

	log.Printf(">>>>>>>>>>>>>>>>> calling consume methods of all workers\n\n")
	workers := []workers.Worker{notificationBatchWorker, schedulerWorker, emailWorker}
	for _, w := range workers {
		go w.Consume()
	}
	log.Printf(">>>>>>>>>>>>>>>>> consume methods called successfully for all workers\n\n")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Workers stopped successfully")
}