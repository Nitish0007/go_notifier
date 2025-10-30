package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nitish0007/go_notifier/internal/workers"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile) // configuring logger to print filename and line number
	log.Println("\n\nStarting Workers...")

	// get all workers
	workers := []func() {
		workers.ConsumeBulkNotificationCreation,
		workers.ConsumeEmailNotifications,
	}

	// start all workers
	for _, worker := range workers {
		go worker()
	}

	// wait for all workers to finish
	<-make(chan bool)

	log.Println("Workers started successfully")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Workers stopped successfully")
}