package contact

import "time"

type ContactResponse struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	AccountID int       `json:"account_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type EmailContactPayload struct {
	Email     string `json:"email" validate:"required,email"`
	ContactID int    `json:"contact_id" validate:"required,gt=0"`
	AccountID int    `json:"account_id" validate:"required,gt=0"`
}

type ContactPayload struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
	AccountID int    `json:"account_id" validate:"required,gt=0"`
	Email     string `json:"email" validate:"required,email"`
	// UUID is optional; when empty, PostgreSQL applies DEFAULT gen_random_uuid() on insert.
	UUID string `json:"uuid,omitempty" validate:"omitempty,uuid"`
}

type CreateContactRequest struct {
	Contact ContactPayload `json:"contact" validate:"required"`
}

type UpdateContactRequest struct {
	Contact ContactPayload `json:"contact" validate:"required"`
}
