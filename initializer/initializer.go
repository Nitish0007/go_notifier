package initializer

import (
	"os"
	"log"
	"os/signal"
	"syscall"
	"context"
	"gorm.io/gorm"
	"github.com/go-chi/chi/v5"
	"github.com/Nitish0007/go_notifier/initializer/container"
	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/emailnotification"
	"github.com/Nitish0007/go_notifier/internal/features/configuration"
	"github.com/Nitish0007/go_notifier/internal/common/middlewares"

	rbmq "github.com/rabbitmq/amqp091-go"
	"github.com/Nitish0007/go_notifier/internal/workers"
)

func InitializeApplication(db *gorm.DB, router *chi.Mux) {
	// Initialize container with all dependencies
	c := container.NewContainer(db)

	// register routes
	router.Route("/api/v1", func(r chi.Router) {
		account.RegisterPublicAccountRoutes(db, r, c.AccountHandler)
		// account.RegisterAccountRoutes(db, r, c.AccountHandler)

		r.Route("/{account_id}", func(r chi.Router) {
			r.Use(middlewares.AuthenticateRequest(db))
			emailnotification.RegisterEmailNotificationRoutes(db, r, c.EmailNotificationHandler)
			configuration.RegisterConfigurationRoutes(db, r, c.ConfigurationHandler)
		})
	})
}

func InitializeWorkers(db *gorm.DB, rbmqConn *rbmq.Connection, ctx context.Context) {
	// Initialize container with all dependencies
	c := container.NewContainer(db)

	// Initialize workers by injecting dependencies
	schedulerWorker := workers.NewNotificationSchedulerWorker(db, rbmqConn, ctx, c.EmailNotificationService)
	emailWorker := workers.NewEmailWorker(db, rbmqConn, ctx, c.EmailNotificationService)
	// notificationBatchWorker := workers.NewNotificationBatchWorker(db, rbmqConn, ctx, c.NotificationService)


	log.Printf(">>>>>>>>>>>>>>>>> calling consume methods of all workers\n\n")
	workers := []workers.Worker{schedulerWorker, emailWorker}
	for _, w := range workers {
		go w.Consume()
	}
	log.Printf(">>>>>>>>>>>>>>>>> consume methods called successfully for all workers\n\n")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Workers stopped successfully")
}