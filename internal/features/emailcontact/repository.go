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

func (r *EmailContactRepository) FindById(ctx context.Context, accId int64, id int64) (*EmailContact, error) {
	var emailContact EmailContact
	err := r.DB.WithContext(ctx).Where("id = ? AND account_id = ?", id, accId).First(&emailContact).Error
	if err != nil {
		return nil, err
	}
	return &emailContact, nil
}

func (r *EmailContactRepository) FindByEmail(ctx context.Context, accId int64, email string) (*EmailContact, error) {
	var emailContact EmailContact
	err := r.DB.WithContext(ctx).Where("account_id = ? AND email = ?", accId, email).First(&emailContact).Error
	if err != nil {
		return nil, err
	}
	return &emailContact, nil
}

// FindByEmailWithTx looks up an email contact using the given DB session or transaction.
func (r *EmailContactRepository) FindByEmailWithTx(ctx context.Context, tx *gorm.DB, accID int64, email string) (*EmailContact, error) {
	var emailContact EmailContact
	err := tx.WithContext(ctx).Where("account_id = ? AND email = ?", accID, email).First(&emailContact).Error
	if err != nil {
		return nil, err
	}
	return &emailContact, nil
}

func (r *EmailContactRepository) GetEmailContactsByContactID(ctx context.Context, accId int64, contactID int64) (*EmailContact, error) {
	var emailContact EmailContact
	err := r.DB.WithContext(ctx).Where("contact_id = ? AND account_id = ?", contactID, accId).First(&emailContact).Error
	if err != nil {
		return nil, err
	}
	return &emailContact, nil
}
