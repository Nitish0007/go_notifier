package services

import (
	"context"
	"github.com/Nitish0007/go_notifier/internal/repositories"
	"github.com/Nitish0007/go_notifier/internal/models"
	// "github.com/Nitish0007/go_notifier/utils"
)

type ConfigurationService struct {
	configurationRepo *repositories.ConfigurationRepository
}

func NewConfigurationService(cr *repositories.ConfigurationRepository) *ConfigurationService {
	return &ConfigurationService{
		configurationRepo: cr,
	}
}

func (s *ConfigurationService) GetConfigurations(ctx context.Context, accID int) ([]*models.Configuration, error) {
	configs, err := s.configurationRepo.Index(ctx, accID)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (s *ConfigurationService) CreateConfiguration(ctx context.Context, configData map[string]any) (*models.Configuration, error) {
	config, err := models.NewConfiguration(configData)
	if err != nil {
		return nil, err
	}
	err = s.configurationRepo.Create(ctx, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// func (s *ConfigurationService) UpdateConfiguration(ctx context.Context, configData map[string]any) (*models.Configuration, error) {
// 	config, err := models.NewConfiguration(configData)
// 	if err != nil {
// 		return nil, err
// 	}
// 	config.ID = configData["id"].(int)
// 	if err := s.configurationRepo.Update(ctx, config); err != nil {
// 		return nil, err
// 	}
// 	return config, nil	
// }