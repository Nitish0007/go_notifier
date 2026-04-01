package contact

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nitish0007/go_notifier/internal/features/emailcontact"
	"github.com/google/uuid"
)

type ContactService struct {
	contactRepository      *ContactRepository
	emailContactRepository *emailcontact.EmailContactRepository
}

func NewContactService(contactRepository *ContactRepository) *ContactService {
	return &ContactService{
		contactRepository: contactRepository,
	}
}

func (s *ContactService) GetContacts(ctx context.Context, accID int) ([]*ContactResponse, error) {
	contacts, err := s.contactRepository.GetContacts(ctx, accID)
	if err != nil {
		return nil, err
	}

	contactResponses := make([]*ContactResponse, len(contacts))
	for i, contact := range contacts {
		contactResponses[i] = &ContactResponse{
			ID:        contact.ID,
			UUID:      contact.UUID,
			AccountID: contact.AccountID,
			FirstName: contact.FirstName,
			LastName:  contact.LastName,
		}
	}
	return contactResponses, nil
}

func (s *ContactService) CreateContact(ctx context.Context, payload *CreateContactRequest) (*ContactResponse, error) {
	contact := &Contact{
		FirstName: payload.Contact.FirstName,
		LastName:  payload.Contact.LastName,
		AccountID: payload.Contact.AccountID,
	}
	if u := strings.TrimSpace(payload.Contact.UUID); u != "" {
		parsed, err := uuid.Parse(u)
		if err != nil {
			return nil, fmt.Errorf("invalid uuid: %w", err)
		}
		contact.UUID = parsed.String()
	}

	ec := &emailcontact.EmailContact{
		Email:     payload.Contact.Email,
		AccountID: payload.Contact.AccountID,
	}

	err := s.contactRepository.CreateWithEmail(ctx, contact, ec)
	if err != nil {
		return nil, err
	}

	return &ContactResponse{
		ID:        contact.ID,
		UUID:      contact.UUID,
		AccountID: contact.AccountID,
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
		Email:     ec.Email,
	}, nil
}

func (s *ContactService) GetContactByKey(ctx context.Context, key string, value any) (*ContactResponse, error) {
	switch key {
	case "id":
		contact, err := s.contactRepository.FindById(ctx, value.(int))
		if err != nil {
			return nil, err
		}
		return &ContactResponse{
			ID:        contact.ID,
			UUID:      contact.UUID,
			AccountID: contact.AccountID,
			FirstName: contact.FirstName,
			LastName:  contact.LastName,
			Email:     contact.EmailContact.Email,
			CreatedAt: contact.CreatedAt,
			UpdatedAt: contact.UpdatedAt,
		}, nil
	case "uuid":
		contact, err := s.contactRepository.FindByUUID(ctx, value.(string))
		if err != nil {
			return nil, err
		}
		return &ContactResponse{
			ID:        contact.ID,
			UUID:      contact.UUID,
			AccountID: contact.AccountID,
			FirstName: contact.FirstName,
			LastName:  contact.LastName,
			Email:     contact.EmailContact.Email,
			CreatedAt: contact.CreatedAt,
			UpdatedAt: contact.UpdatedAt,
		}, nil
	default:
		return nil, fmt.Errorf("invalid key: %s", key)
	}
}
