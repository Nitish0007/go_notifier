package models

import (
	"time"
)

type ApiKey struct {
	ID        	int       	`json:"id"`
	Key       	string    	`json:"key"`
	AccountID 	int       	`json:"account_id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}