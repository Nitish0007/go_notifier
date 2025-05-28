package main

import (
	"log"
	"net/http"

	"github.com/Nitish0007/go_notifier/utils"
	"github.com/Nitish0007/go_notifier/initializer"
)

func main(){
	log.SetFlags(log.LstdFlags | log.Llongfile) // configuring logger to print filename and line number

	conn, _ := utils.ConnectDB()
	r := utils.InitRouter()
	initializer.InititalizeApplication(conn, r)

	http.ListenAndServe(":8080", r)
}

// command to create migration
// -> migrate create -ext sql -dir db/migrations -seq migration_name