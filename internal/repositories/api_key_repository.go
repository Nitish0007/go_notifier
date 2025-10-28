package repositories

import (
	"context"
	"gorm.io/gorm"

	"github.com/Nitish0007/go_notifier/internal/models"
)


type ApiKeyRepository struct {
	DB *gorm.DB
}

func NewApiKeyRepository(conn *gorm.DB) *ApiKeyRepository{
	return &ApiKeyRepository{
		DB: conn,
	}
}

func (r *ApiKeyRepository) Create(ctx context.Context, apiKey *models.ApiKey) error {
	err := r.DB.WithContext(ctx).Create(apiKey).Error; if err != nil {
		return err
	}
	return nil
}

func (r *ApiKeyRepository) FindByAccountID(ctx context.Context, accountID int) (*models.ApiKey, error) {
	var apiKey models.ApiKey
	err := r.DB.WithContext(ctx).Where("account_id = ?", accountID).First(&apiKey).Error; if err != nil {
		return nil, err
	}
	return &apiKey, nil
}