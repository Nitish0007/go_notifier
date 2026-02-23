package repositories

import (
	"context"
	"fmt"
	"github.com/Nitish0007/go_notifier/internal/models"
	"gorm.io/gorm"
)

type ConfigurationRepository struct {
	DB *gorm.DB
}

func NewConfigurationRepository(conn *gorm.DB) *ConfigurationRepository {
	return &ConfigurationRepository{
		DB: conn,
	}
}

func (r *ConfigurationRepository) Create(ctx context.Context, config *models.Configuration) error {
	return r.DB.WithContext(ctx).Create(config).Error
}

func (r *ConfigurationRepository) GetByAccountID(ctx context.Context, accountID int) (*models.Configuration, error) {
	var config models.Configuration
	err := r.DB.WithContext(ctx).Where("account_id = ? AND default_configuration = ?", accountID, true).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ConfigurationRepository) GetByAccountIDAndConfigType(ctx context.Context, accountID int, configType string) (*models.Configuration, error) {
	var config models.Configuration
	err := r.DB.WithContext(ctx).Where("account_id = ? AND config_type = ? AND default_configuration = ?", accountID, configType, false).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ConfigurationRepository) Index(ctx context.Context, accountID int) ([]*models.Configuration, error) {
	var configs []*models.Configuration
	err := r.DB.WithContext(ctx).Where("account_id = ?", accountID).Order("created_at DESC").Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *ConfigurationRepository) Update(ctx context.Context, config *models.Configuration) error {
	var existingConfig models.Configuration
	err := r.DB.WithContext(ctx).Where("id = ? AND account_id = ?", config.ID, config.AccountID).First(&existingConfig).Error
	if err != nil {
		return err
	}
	
	existingConfig.ConfigType = config.ConfigType
	existingConfig.ConfigurationData = config.ConfigurationData
	existingConfig.DefaultConfiguration = config.DefaultConfiguration

	err = r.DB.WithContext(ctx).Save(&existingConfig).Error
	if err != nil {
		return err
	}
	config = &existingConfig
	return nil
}

func (r *ConfigurationRepository) Delete(ctx context.Context, id int) error {
	result := r.DB.WithContext(ctx).Where("id = ?", id).Delete(&models.Configuration{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("configuration not found with id: %d", id)
	}
	return nil
}

func (r *ConfigurationRepository) GetByFields(ctx context.Context, fields map[string]any) (*models.Configuration, error) {
	var config models.Configuration
	err := r.DB.WithContext(ctx).Where(fields).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}