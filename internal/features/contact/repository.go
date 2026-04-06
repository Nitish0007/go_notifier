package contact

import (
	"fmt"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Nitish0007/go_notifier/internal/features/emailcontact"
)

type ContactRepository struct {
	DB               *gorm.DB
	emailContactRepo *emailcontact.EmailContactRepository
}

func NewContactRepository(db *gorm.DB, ecRepo *emailcontact.EmailContactRepository) *ContactRepository {
	return &ContactRepository{
		DB:               db,
		emailContactRepo: ecRepo,
	}
}

func (r *ContactRepository) GetContacts(ctx context.Context, accID int64) ([]*Contact, error) {
	var contacts []*Contact
	err := r.DB.WithContext(ctx).Preload("EmailContact").Where("account_id = ?", accID).Find(&contacts).Error
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (r *ContactRepository) CreateWithEmail(ctx context.Context, contact *Contact, ec *emailcontact.EmailContact) error {
	result := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.createContactRow(tx, contact); err != nil {
			return err
		}

		var existingEmailContact = &emailcontact.EmailContact{}
		if err := tx.Where("email = ? AND account_id = ?", ec.Email, ec.AccountID).
			First(&existingEmailContact).Error; err == nil {
			// Email already exists for this account, skip inserting or return application error.
			return errors.New("email already exists for this account: " + ec.Email)
		}

		// assign contact ID
		ec.ContactID = contact.ID

		return tx.Create(ec).Error
	})
	return result
}

// createContactRow inserts contact. If UUID is empty, the column is omitted on PostgreSQL so
// DEFAULT gen_random_uuid() applies; SQLite has no such default in tests, so a UUID is generated in-app.
func (r *ContactRepository) createContactRow(tx *gorm.DB, contact *Contact) error {
	if contact.UUID != "" {
		return tx.Create(contact).Error
	}
	if tx.Dialector.Name() == "sqlite" {
		contact.UUID = uuid.New().String()
		return tx.Create(contact).Error
	}
	if err := tx.Omit("UUID").Create(contact).Error; err != nil {
		return err
	}
	return tx.Where("id = ?", contact.ID).First(contact).Error
}

func (r *ContactRepository) FindById(ctx context.Context, accId int64, id int64) (*Contact, error) {
	var contact Contact
	err := r.DB.WithContext(ctx).Preload("EmailContact").Where("id = ? AND account_id = ?", id, accId).First(&contact).Error
	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func (r *ContactRepository) FindByUUID(ctx context.Context, accId int64, uuid string) (*Contact, error) {
	var contact Contact
	err := r.DB.WithContext(ctx).Preload("EmailContact").Where("account_id = ? AND uuid = ?", accId, uuid).First(&contact).Error
	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func (r *ContactRepository) FindByEmail(ctx context.Context, accId int64, email string) (*Contact, error) {
	var contact Contact
	err := r.DB.WithContext(ctx).Where("account_id = ? AND email_contact->>'email' = ?", accId, email).First(&contact).Error
	if err != nil {
		return nil, err
	}
	if contact.EmailContact == nil {
		return nil, nil
	}
	return &contact, nil
}

func (r *ContactRepository) FindOrCreateByEmail(ctx context.Context, accId int64, contactPayload *ContactPayload) (*Contact, error) {
	contact, err := r.FindByEmail(ctx, accId, contactPayload.Email)
	if err != nil {
		return nil, err
	}
	if contact == nil {
		contact = &Contact{
			FirstName: contactPayload.FirstName,
			LastName: contactPayload.LastName,
			AccountID: accId,
		}
		emailContact := &emailcontact.EmailContact{
			Email:     contactPayload.Email,
			AccountID: accId,
		}
		err = r.CreateWithEmail(ctx, contact, emailContact)
		if err != nil {
			return nil, fmt.Errorf("failed to create contact: %w", err)
		}
	}
	return contact, nil
}