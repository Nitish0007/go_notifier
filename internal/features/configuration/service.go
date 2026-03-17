package configuration

import (
	"context"
	"log"
)

type ConfigurationService struct {
	configurationRepo *ConfigurationRepository
}

func NewConfigurationService(cr *ConfigurationRepository) *ConfigurationService {
	return &ConfigurationService{
		configurationRepo: cr,
	}
}

func (s *ConfigurationService) GetConfigurations(ctx context.Context, accID int) ([]*Configuration, error) {
	configs, err := s.configurationRepo.Index(ctx, accID)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (s *ConfigurationService) CreateConfiguration(ctx context.Context, payload *CreateConfigurationRequest) (*Configuration, error) {
	configRequest := ConfigurationRequest{
		Configuration: ConfigurationPayload{
			AccountID: payload.Configuration.AccountID,
			DefaultConfiguration: payload.Configuration.DefaultConfiguration,
			ConfigType: payload.Configuration.ConfigType,
			ConfigurationData: payload.Configuration.ConfigurationData,
		},
	}
	config, err := NewConfiguration(&configRequest)
	if err != nil {
		return nil, err
	}
	err = s.configurationRepo.Create(ctx, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (s *ConfigurationService) DeleteConfiguration(ctx context.Context, accID int, cid int) error {
	return s.configurationRepo.Delete(ctx, cid)
}

func (s *ConfigurationService) UpdateConfiguration(ctx context.Context, payload *UpdateConfigurationRequest) (*Configuration, error) {
	configRequest := ConfigurationRequest{
		Configuration: ConfigurationPayload{
			ID: payload.Configuration.ID,
			AccountID: payload.Configuration.AccountID,
			DefaultConfiguration: payload.Configuration.DefaultConfiguration,
			ConfigType: payload.Configuration.ConfigType,
			ConfigurationData: payload.Configuration.ConfigurationData,
		},
	}
	log.Printf(">>>>>>>>>>>>>> Service - config: %+v", payload.Configuration)
	config, err := NewConfiguration(&configRequest)
	if err != nil {
		return nil, err
	}
	err = s.configurationRepo.Update(ctx, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}