package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	// Initialize Services by injecting corresponding repository dependency
	blkNotificationService := services.NewBulkNotificationService(notificationRepo, notificationBatchRepo, notificationBatchErrorRepo)

	// Initialize workers by injecting dependencies
	notificationBatchWorker := workers.NewNotificationBatchWorker(dbConn, rbmqConn, ctx, blkNotificationService)

	// start workers by calling Consume method
	notificationBatchWorker.Consume()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Workers stopped successfully")
}