package initializer

import (
	"os"
	"log"
	"os/signal"
	"syscall"
	"context"
	"gorm.io/gorm"
	"google.golang.org/grpc"
	"github.com/go-chi/chi/v5"
	"github.com/Nitish0007/go_notifier/initializer/container"
	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/emailnotification"
	"github.com/Nitish0007/go_notifier/internal/features/content"
	"github.com/Nitish0007/go_notifier/internal/features/list"
	"github.com/Nitish0007/go_notifier/internal/features/contact"
	"github.com/Nitish0007/go_notifier/internal/features/configuration"
	"github.com/Nitish0007/go_notifier/internal/common/middlewares"
	"github.com/Nitish0007/go_notifier/internal/common/mq"
	"github.com/Nitish0007/go_notifier/internal/workers"
	accountv1 "github.com/Nitish0007/go_notifier/pkg/gen/account/v1"
	"github.com/Nitish0007/go_notifier/internal/common/rabbitmq"
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
      list.RegisterListRoutes(db, r, c.ListHandler)
			contact.RegisterContactRoutes(db, r, c.ContactHandler)
			content.RegisterContentRoutes(db, r, c.ContentHandler)
		})
	})
}

func InitializeGRPCServer(db *gorm.DB, grpcServer *grpc.Server) {
	c := container.NewContainer(db)
	grpcAccount := account.NewAccountServiceServer(c.AccountService)
	// register gRPC handlers
	accountv1.RegisterAccountServiceServer(grpcServer, grpcAccount)
}

func InitializeWorkers(db *gorm.DB, mqClient mq.MQClient, ctx context.Context) {
	// Initialize container with all dependencies
	c := container.NewContainer(db)

	// initialize queues
	_, err := rabbitmq.InitializeQueues(mqClient)
	if err != nil {
		log.Fatalf("Failed to initialize queues: %v", err)
	}

	// Initialize workers by injecting dependencies
	emailSchedulerWorker := workers.NewEmailSchedulerWorker(ctx, db, mqClient)
	// schedulerWorker := workers.NewNotificationSchedulerWorker(db, mqClient, ctx, c.EmailNotificationService)
	emailWorker := workers.NewEmailWorker(db, mqClient, ctx, c.EmailNotificationService)
	// notificationBatchWorker := workers.NewNotificationBatchWorker(db, rbmqConn, ctx, c.NotificationService)


	log.Printf(">>>>>>>>>>>>>>>>> calling consume methods of all workers\n\n")
	workers := []workers.Worker{emailSchedulerWorker, emailWorker}
	for _, w := range workers {
		go w.Run()
	}
	log.Printf(">>>>>>>>>>>>>>>>> consume methods called successfully for all workers\n\n")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Workers stopped successfully")
}