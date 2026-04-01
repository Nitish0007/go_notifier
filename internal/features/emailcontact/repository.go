package emailcontact

import (
	"context"
	"gorm.io/gorm"
)

type EmailContactRepository struct {
	DB *gorm.DB
}

func NewEmailContactRepository(db *gorm.DB) *EmailContactRepository {
	return &EmailContactRepository{
		DB: db,
	}
}

func (r *EmailContactRepository) FindById(ctx context.Context, id int) (*EmailContact, error) {
	var emailContact EmailContact
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&emailContact).Error
	if err != nil {
		return nil, err
	}
	return &emailContact, nil
}

func (r *EmailContactRepository) FindByEmail(ctx context.Context, email string) (*EmailContact, error) {
	var emailContact EmailContact
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&emailContact).Error
	if err != nil {
		return nil, err
	}
	return &emailContact, nil
}

func (r *EmailContactRepository) GetEmailContactsByContactID(ctx context.Context, contactID int) (*EmailContact, error) {
	var emailContact EmailContact
	err := r.DB.WithContext(ctx).Where("contact_id = ?", contactID).First(&emailContact).Error
	if err != nil {
		return nil, err
	}
	return &emailContact, nil
}
