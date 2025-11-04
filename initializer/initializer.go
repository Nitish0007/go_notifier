package initializer

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	// "github.com/jackc/pgx/v5/pgxpool"

	"github.com/Nitish0007/go_notifier/internal/handlers"
	"github.com/Nitish0007/go_notifier/internal/middlewares"
	"github.com/Nitish0007/go_notifier/internal/notifiers"
	"github.com/Nitish0007/go_notifier/internal/repositories"
	"github.com/Nitish0007/go_notifier/internal/routes"
	"github.com/Nitish0007/go_notifier/internal/services"
)

func InititalizeApplication(db *gorm.DB, router *chi.Mux) {
	// initializing repositories, services and handlers by injecting dependencies

	// Intialize Repositories by injecting db connection dependency
	accRepo := repositories.NewAccountRepository(db)
	apiKeyRepo := repositories.NewApiKeyRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	notificationBatchRepo := repositories.NewNotificationBatchRepository(db)
	notificationBatchErrorRepo := repositories.NewNotificationBatchErrorRepo(db)

	// intialize notifiers
	emailNotifier := notifiers.NewEmailNotifier(notificationRepo)

	// Initialize Services by injecting corresponding repository dependency
	accService := services.NewAccountService(accRepo, apiKeyRepo)
	notificationService := services.NewNotificationService(
		[]notifiers.Notifier{emailNotifier},
		notificationRepo,
	)
	bulkNotificationService := services.NewBulkNotificationService(notificationRepo, notificationBatchRepo, notificationBatchErrorRepo)

	// Initialize Handlers by injecting corresponding service dependency
	accountHandler := handlers.NewAccountHandler(accService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	bulkNotificationHandler := handlers.NewBulkNotificationHandler(bulkNotificationService)
	// Register Routes by injecting corresponding handler dependency
	router.Route("/api/v1", func(r chi.Router) {
		routes.RegisterPublicAccountRoutes(db, r, accountHandler)

		// protected routes
		r.Route("/{account_id}", func(authenticated chi.Router) {
			authenticated.Use(middlewares.AuthenticateRequest(db))

			routes.RegisterAccountRoutes(db, authenticated, accountHandler)
			routes.RegisterNotificationRoutes(db, authenticated, notificationHandler)
			routes.RegisterBulkNotificationRoutes(db, authenticated, bulkNotificationHandler)
		})
	})
}
