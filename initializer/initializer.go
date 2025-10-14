package initializer

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Nitish0007/go_notifier/internal/handlers"
	"github.com/Nitish0007/go_notifier/internal/middlewares"
	"github.com/Nitish0007/go_notifier/internal/notifiers"
	"github.com/Nitish0007/go_notifier/internal/repositories"
	"github.com/Nitish0007/go_notifier/internal/routes"
	"github.com/Nitish0007/go_notifier/internal/services"
)

func InititalizeApplication(connPool *pgxpool.Pool, router *chi.Mux) {
	// initializing repositories, services and handlers by injecting dependencies

	// Intialize Repositories by injecting db connection dependency
	accRepo := repositories.NewAccountRepository(connPool)
	apiKeyRepo := repositories.NewApiKeyRepository(connPool)
	notificationRepo := repositories.NewNotificationRepository(connPool)

	// intialize notifiers
	emailNotifier := notifiers.NewEmailNotifier(notificationRepo)

	// Initialize Services by injecting corresponding repository dependency
	accService := services.NewAccountService(accRepo, apiKeyRepo)
	notificationService := services.NewNotificationService(
		[]notifiers.Notifier{emailNotifier},
		notificationRepo,
	)
	bulkNotificationService := services.NewBulkNotificationService(notificationRepo)

	// Initialize Handlers by injecting corresponding service dependency
	accountHandler := handlers.NewAccountHandler(accService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	bulkNotificationHandler := handlers.NewBulkNotificationHandler(bulkNotificationService)
	// Register Routes by injecting corresponding handler dependency
	router.Route("/api/v1", func(r chi.Router) {
		routes.RegisterPublicAccountRoutes(connPool, r, accountHandler)

		// protected routes
		r.Route("/{account_id}", func(authenticated chi.Router) {
			authenticated.Use(middlewares.AuthenticateRequest(connPool))

			routes.RegisterAccountRoutes(connPool, authenticated, accountHandler)
			routes.RegisterNotificationRoutes(connPool, authenticated, notificationHandler)
			routes.RegisterBulkNotificationRoutes(connPool, authenticated, bulkNotificationHandler)
		})
	})
}
