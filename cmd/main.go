package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Nitish0007/go_notifier/initializer"
	"github.com/Nitish0007/go_notifier/utils"
)

func main(){
	log.SetFlags(log.LstdFlags | log.Llongfile) // configuring logger to print filename and line number

	conn, _ := utils.ConnectDB()
	defer conn.Close(context.Background())
	
	r := utils.InitRouter()
	initializer.InititalizeApplication(conn, r)

	http.ListenAndServe(":8080", r)
}

// command to create migration
// -> migrate create -ext sql -dir db/migrations -seq migration_name