package list

import "time"

// Request DTOs
type CreateListRequest struct {
	List struct {
		AccountID int64  `json:"account_id" validate:"required,gt=0"`
		Name      string `json:"name" validate:"required,min=1,max=100"`
	} `json:"list" validate:"required"`
}

type SubscribeToListRawPayload struct {
	ContactID  int64  `json:"contact_id" validate:"omitempty,gt=0"` // if contact already exists, user can send contact_id or uuid in that case
	UUID string `json:"uuid" validate:"omitempty,uuid"`
	Active bool `json:"active" default:"false" validate:"omitempty,boolean"`
  EmailContact struct {
		UUID string `json:"uuid" validate:"omitempty,uuid"`
		FirstName string `json:"first_name" validate:"omitempty,min=1,max=100"`
		LastName string `json:"last_name" validate:"omitempty,min=1,max=100"`
		Email string `json:"email" validate:"omitempty,email"`
	} `json:"email_contact" validate:"omitempty"` // this is to provide ability to create a new contact in the system and add it to list as well
}

type SubscribeToListRequest struct {
	AccountID  int64      `json:"account_id" validate:"required, gt=0"`
	ListID     int64      `json:"list_id" validate:"required,gt=0"`
	ContactID  int64      `json:"contact_id" validate:"required,gt=0"`
	Active     bool       `json:"active" validate:"required,boolean"`
	CreatedAt  time.Time  `json:"created_at" validate:"-"`
	UpdatedAt  time.Time  `json:"updated_at" validate:"-"`
}

// Response DTOs
type CreateListResponse struct {
	ID            int64     `json:"id"`
	AccountID     int64     `json:"account_id"`
	Name          string    `json:"name"`
	ContactsCount int64     `json:"contacts_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SubscribeToListResponse struct {
	ListID int64 `json:"list_id"`
	ContactID  int64  `json:"contact_id"`
	ContactUUID string `json:"contact_uuid"`
	Active bool `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubscriberResponse struct {
	FirstName             string     `json:"first_name"`
	LastName              string     `json:"last_name"`
	ContactID             int64      `json:"contact_id"`
	UUID                  string     `json:"uuid"`
	// Email                 string     `json:"email"`
	Active                bool       `json:"active"`
	ContactCreatedAt      time.Time  `json:"contact_created_at"`
	ContactUpdatedAt      time.Time  `json:"contact_updated_at"`
	SubscriptionCreatedAt time.Time  `json:"subscription_created_at"`
	SubscriptionUpdatedAt time.Time  `json:"subscription_updated_at"`
}

type ListOfSubscribersResponse struct {
	ListID      int64                  `json:"list_id"`
	AccountID   int64                  `json:"account_id"`
	Subscribers []*SubscriberResponse   `json:"subscribers"`
}