package api_key

import (
	"context"
	"gorm.io/gorm"
)

type ApiKeyRepository struct {
	DB *gorm.DB
}

func NewApiKeyRepository(conn *gorm.DB) *ApiKeyRepository{
	return &ApiKeyRepository{
		DB: conn,
	}
}

func (r *ApiKeyRepository) Create(ctx context.Context, apiKey *ApiKey) error {
	err := r.DB.WithContext(ctx).Create(apiKey).Error; if err != nil {
		return err
	}
	return nil
}

func (r *ApiKeyRepository) FindByAccountID(ctx context.Context, accountID int) (ApiKey, error) {
	var apiKey ApiKey
	err := r.DB.WithContext(ctx).Where("account_id = ?", accountID).First(&apiKey).Error; if err != nil {
		return ApiKey{}, err
	}
	return apiKey, nil
}