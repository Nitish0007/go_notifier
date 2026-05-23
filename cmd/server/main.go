package main

import (
	"os"
	"net"
	"log"
	"sync"
	"time"
	"syscall"
	"context"
	"net/http"
	"os/signal"
	"google.golang.org/grpc"
	"github.com/joho/godotenv"
	"github.com/go-chi/chi/v5"
	"github.com/Nitish0007/go_notifier/initializer"
	"github.com/Nitish0007/go_notifier/internal/common/router"
	"github.com/Nitish0007/go_notifier/internal/common/database"
	"github.com/Nitish0007/go_notifier/internal/common/interceptors"
)

// For printing all registered routes - helpful in debugging
func PrintRoutes(r chi.Router) {
	err := chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking routes: %s\n", err)
	}
}

func main() {
  ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	
	log.Println("\n\nStarting API Server...")
	log.SetFlags(log.LstdFlags | log.Llongfile) // configuring logger to print filename and line number

	// load environment variables
	env := os.Getenv("ENV")
	
	if env == "development" {
		err := godotenv.Load(".env.development")
		if err != nil {
			log.Fatalf("Error loading .env.development file: %v", err)
		}
	}

	db, err := database.Connect(env)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	r := router.InitRouter()
	initializer.InitializeApplication(db, r)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.UnaryInterceptor,
			interceptors.AuthUnaryInterceptor(db),
		),
	)
	initializer.InitializeGRPCServer(db, grpcServer)
	// PrintRoutes(r)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("HTTP server starting...")
		httpServer.ListenAndServe()
	}()
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("gRPC server starting...")
		grpcServer.Serve(listener)
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received, shutting down servers...")

	// HTTP shutdown requires a context with a timeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP shutdown error: %v", err)
	}

	// gRPC provides a built-in GracefulStop.
	grpcServer.GracefulStop()

	wg.Wait()
	log.Println("Servers exited gracefully.")
}

// command to create migration
// -> migrate create -ext sql -dir db/migrations -seq migration_name
