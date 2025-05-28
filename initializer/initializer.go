package initializer

import (
	"github.com/jackc/pgx/v5"
	"github.com/go-chi/chi/v5"
	
	"github.com/Nitish0007/go_notifier/internal/routes"
	"github.com/Nitish0007/go_notifier/internal/services"
	"github.com/Nitish0007/go_notifier/internal/handlers"
	"github.com/Nitish0007/go_notifier/internal/repositories"
)

func InititalizeApplication(conn *pgx.Conn, router *chi.Mux){
	// initializing repositories, services and handlers by injecting dependencies

	// Intialize Repositories by injecting db connection dependency
	accRepo := repositories.NewAccountRepository(conn)
	apiKeyRepo := repositories.NewApiKeyRepository(conn)

	// Initialize Services by injecting corresponding repository dependency
	accService := services.NewAccountService(accRepo, apiKeyRepo)

	// Initialize Handlers by injecting corresponding service dependency
	accountHandler := handlers.NewAccountHandler(accService)

	// Register Routes by injecting corresponding handler dependency
	routes.RegisterAccountRoutes(router, accountHandler)

}