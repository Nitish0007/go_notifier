package account

import (
	"context"

	"gorm.io/gorm"

	"github.com/Nitish0007/go_notifier/utils"
	"github.com/Nitish0007/go_notifier/internal/features/apiKey"
)

type AccountRepository struct {
	DB               *gorm.DB
	ApiKeyRepository *apiKey.ApiKeyRepository
}

func NewAccountRepository(conn *gorm.DB, apiKeyRepository *apiKey.ApiKeyRepository) *AccountRepository {
	return &AccountRepository{
		DB:               conn,
		ApiKeyRepository: apiKeyRepository,
	}
}

// creates a new account without API key
func (r *AccountRepository) Create(ctx context.Context, account *Account) error {
	result := r.DB.WithContext(ctx).Create(account)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// creates a new account and API key within a transaction as a atomic operation.
func (r *AccountRepository) RegisterAccount(ctx context.Context, account *Account) error {
	result := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(account).Error; err != nil {
			return err
		}

		// initialize API key
		apiKey := &apiKey.ApiKey{
			Key:       utils.GenerateAlphaNumericKey(),
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

func (r *AccountRepository) GetApiKeyByAccountID(ctx context.Context, accountID int) (string, error) {
	var apiKey apiKey.ApiKey
	result := r.DB.WithContext(ctx).Where("account_id = ?", accountID).First(&apiKey)
	if result.Error != nil {
		return "", result.Error
	}
	return apiKey.Key, nil
}
func (r *AccountRepository) FindAccountByEmail(ctx context.Context, email string) (Account, error) {
	var account Account
	result := r.DB.WithContext(ctx).Where("email = ?", email).First(&account)
	if result.Error != nil {
		return Account{}, result.Error
	}
	return account, nil
}
