package content

import "time"

// Request DTOs
type CreateContentRequest struct {
	Content struct {
		AccountID int64  `json:"account_id" validate:"required,gt=0"`
		Name string `json:"name" validate:"required"`
		Body      string `json:"body" validate:"required"`
	} `json:"content" validate:"required"`
}

// Response DTOs
type CreateContentResponse struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	AccountID int64 `json:"account_id"`
	Body string `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}