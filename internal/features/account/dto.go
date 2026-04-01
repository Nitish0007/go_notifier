package account

import "time"

// Request DTOs
type SignupRequest struct {
	Account struct {
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required,min=8"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
		FirstName       string `json:"first_name" validate:"required,min=1,max=100"`
		LastName        string `json:"last_name" validate:"required,min=1,max=100"`
	} `json:"account" validate:"required"`
}

type LoginRequest struct {
	Login struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	} `json:"login" validate:"required"`
}

// Response DTOs
type SignupResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginResponse struct {
	AuthToken string `json:"auth_token"`
	AccountID string `json:"account_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
