package main

import (
	"net/http"
	"github.com/Nitish0007/go_notifier/utils"
)

func main(){
	utils.ConnectDB()
	r := utils.InitRouter()

	http.ListenAndServe(":8080", r)
}

// command to create migration
// -> migrate create -ext sql -dir db/migrations -seq migration_name