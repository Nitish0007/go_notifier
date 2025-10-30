package main

import (
	// "context"
	"log"
	"net/http"

	"github.com/Nitish0007/go_notifier/initializer"
	"github.com/Nitish0007/go_notifier/utils"
	"github.com/go-chi/chi/v5"
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
	log.Println("\n\nStarting API Server...")
	log.SetFlags(log.LstdFlags | log.Llongfile) // configuring logger to print filename and line number

	db, err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	r := utils.InitRouter()
	initializer.InititalizeApplication(db, r)

	// PrintRoutes(r)

	http.ListenAndServe(":8080", r)
}

// command to create migration
// -> migrate create -ext sql -dir db/migrations -seq migration_name
