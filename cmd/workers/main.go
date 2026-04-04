package main

import (
	"context"
	"log"
	"os"
	"github.com/joho/godotenv"

	"github.com/Nitish0007/go_notifier/initializer"
	"github.com/Nitish0007/go_notifier/internal/common/database"
	rabbitmq_utils "github.com/Nitish0007/go_notifier/utils/rabbitmq"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile) // configuring logger to print filename and line number
	log.Println("\n\nStarting Workers...")

	// load environment variables
	env := os.Getenv("ENV")
	
	if env == "development" {
		err := godotenv.Load(".env.development")
		if err != nil {
			log.Fatalf("Error loading .env.development file: %v", err)
		}
	} else {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	// make database connection for workers
	dbConn, err := database.Connect(env)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// make RabbitMQ connection for workers
	rbmqConn := rabbitmq_utils.ConnectMQ()
	defer rbmqConn.Close()

	// create context for workers
	ctx := context.Background()

	initializer.InitializeWorkers(dbConn, rbmqConn, ctx)
}