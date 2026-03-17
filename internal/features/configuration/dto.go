package configuration

import (
	"time"
)

// Request DTOs
type CreateConfigurationRequest struct {
	Configuration struct {
		AccountID 						int `json:"account_id" validate:"required,gt=0"`
		DefaultConfiguration 	*bool `json:"default_configuration" validate:"required"`
		ConfigType 						string `json:"config_type" validate:"required,oneof=smtp in_app"`
		ConfigurationData 		map[string]any `json:"configuration_data" validate:"required"`
	} `json:"configuration" validate:"required"`
}

type UpdateConfigurationRequest struct {
	Configuration struct {
		ID 											int `json:"id" validate:"omitempty,gt=0"`
		AccountID 							int `json:"account_id" validate:"omitempty,gt=0"`
		DefaultConfiguration 		*bool `json:"default_configuration" validate:"omitempty,boolean"`
		ConfigType 							string `json:"config_type" validate:"omitempty,oneof=smtp in_app"`
		ConfigurationData 			map[string]any `json:"configuration_data" validate:"omitempty"`
	} `json:"configuration" validate:"required"`
}

type ConfigurationPayload struct {
	ID 										int `json:"id" validate:"omitempty,gt=0"`
	AccountID 						int `json:"account_id" validate:"required,gt=0"`
	DefaultConfiguration 	*bool `json:"default_configuration" validate:"omitempty,boolean"`
	ConfigType 						string `json:"config_type" validate:"omitempty,oneof=smtp in_app"`
	ConfigurationData 		map[string]any `json:"configuration_data" validate:"omitempty"`
}

type ConfigurationRequest struct {
	Configuration ConfigurationPayload `json:"configuration" validate:"required"`
}

// Response DTOs
type CreateConfigurationResponse struct {
	ID 										int `json:"id"`
	AccountID 						int `json:"account_id"`
	DefaultConfiguration 	bool `json:"default_configuration"`
	ConfigType 						string `json:"config_type"`
	ConfigurationData 		map[string]any `json:"configuration_data"`
	CreatedAt 						time.Time `json:"created_at"`
	UpdatedAt 						time.Time `json:"updated_at"`
}