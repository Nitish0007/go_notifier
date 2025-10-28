package repositories

import (
	"context"

	"github.com/Nitish0007/go_notifier/internal/models"

	"gorm.io/gorm"
	"github.com/Nitish0007/go_notifier/utils"
)

type AccountRepository struct {
	DB *gorm.DB
}

func NewAccountRepository(conn *gorm.DB) *AccountRepository {
	return &AccountRepository{
		DB: conn,
	}
}

func (r *AccountRepository) Create(ctx context.Context, account *models.Account) error {
	result := r.DB.WithContext(ctx).Create(account)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// creates a new account within a transaction.
func (r *AccountRepository) CreateAccountWithAPIKeyWithinTx(ctx context.Context, account *models.Account, db *gorm.DB) error {
	result := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(account).Error; err != nil {
			return err
		}
		// initialize API key
		apiKey := &models.ApiKey{
			Key:      utils.GenerateAlphaNumericKey(),
			AccountID: account.ID,
		}
		// create API key
		if err := tx.Create(apiKey).Error; err != nil {
			return err
		}
		return nil
	})
	return result
}

func (r *AccountRepository) FindAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	var account models.Account
	result := r.DB.WithContext(ctx).Where("email = ?", email).First(&account)
	if result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}
