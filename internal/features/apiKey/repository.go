package apiKey

import (
	"context"
  "errors"
	"gorm.io/gorm"
)

type ApiKeyRepository struct {
	DB *gorm.DB
}

func NewApiKeyRepository(conn *gorm.DB) *ApiKeyRepository {
	return &ApiKeyRepository{
		DB: conn,
	}
}

func (r *ApiKeyRepository) Index(ctx context.Context, accountID int64) ([]ApiKey, error) {
	var apiKeys []ApiKey
	err := r.DB.WithContext(ctx).Where("account_id = ?", accountID).Order("created_at DESC").Find(&apiKeys).Error
	if err != nil {
		return nil, errors.New("failed to find api keys by account ID: " + err.Error())
	}
	return apiKeys, nil
}

func (r *ApiKeyRepository) Create(ctx context.Context, apiKey *ApiKey) error {
	err := r.DB.WithContext(ctx).Create(apiKey).Error
	if err != nil {
		return errors.New("failed to create api key: " + err.Error())
	}
	return nil
}

func (r *ApiKeyRepository) FindByKeyAndAccountID(ctx context.Context, key string, accountID int64) (ApiKey, error) {
	var apiKey ApiKey
	err := r.DB.WithContext(ctx).Where("key = ? AND account_id = ?", key, accountID).First(&apiKey).Error
	if err != nil {
		return ApiKey{}, errors.New("failed to find api key by: " + err.Error())
	}
	return apiKey, nil
}
