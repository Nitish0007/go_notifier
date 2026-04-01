package configuration

import (
	// "log"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ConfigurationType string

const (
	SMTPConfig   ConfigurationType = "smtp"
	WebAppConfig ConfigurationType = "in_app" // this is for web app notifications
)

type ConfigData struct {
	SMTPConfiguration   SMTPConfiguration   `json:"smtp_configuration" validate:"omitempty"`
	WebAppConfiguration WebAppConfiguration `json:"web_app_configuration" validate:"omitempty"`
}

type SMTPConfiguration struct {
	Host     string `json:"host" validate:"required,min=1"`
	Port     int    `json:"port" validate:"required,gt=0,lte=65535"`
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
	From     string `json:"from" validate:"required,email"`
}

type WebAppConfiguration struct {
	WebAppURL    string `json:"web_app_url" validate:"required,url"`
	WebAppSecret string `json:"web_app_secret" validate:"required,min=1"`
	WebAppToken  string `json:"web_app_token" validate:"required,min=1"`
}

// ToMap returns a flat map for storage/API. Add new config types here with their own ToMap.
func (s SMTPConfiguration) ToMap() map[string]any {
	out := make(map[string]any)
	jsonBytes, _ := json.Marshal(s)
	_ = json.Unmarshal(jsonBytes, &out)
	return out
}

func (s WebAppConfiguration) ToMap() map[string]any {
	out := make(map[string]any)
	jsonBytes, _ := json.Marshal(s)
	_ = json.Unmarshal(jsonBytes, &out)
	return out
}

type Configuration struct {
	ID                   int            `json:"id" gorm:"primaryKey" validate:"-"`
	AccountID            int            `json:"account_id" gorm:"not null;index" validate:"required,gt=0"`
	DefaultConfiguration bool           `json:"default_configuration" gorm:"default:false" validate:"-"`
	ConfigurationData    map[string]any `json:"configuration_data" gorm:"serializer:json" validate:"required"`
	ConfigType           string         `json:"config_type" gorm:"not null" validate:"required,oneof=smtp in_app"`
	CreatedAt            time.Time      `json:"created_at" gorm:"autoCreateTime" validate:"-"`
	UpdatedAt            time.Time      `json:"updated_at" gorm:"autoUpdateTime" validate:"-"`
}

func (c *Configuration) BeforeSave(tx *gorm.DB) error {
	// check uniqueness for account_id, config_type and default_configuration
	if c.DefaultConfiguration {
		var existingConfig Configuration
		err := tx.Where("account_id = ? AND config_type = ? AND default_configuration = ?", c.AccountID, c.ConfigType, true).First(&existingConfig).Error
		if err != nil {
			return err
		}
		if existingConfig.ID != c.ID {
			return errors.New("default configuration already exists")
		}
	}
	return nil
}

// NewConfiguration builds a Configuration from the request. configuration_data must already
// be validated flat by config_type (e.g. via ValidateConfigurationDataByType).
func NewConfiguration(payload *ConfigurationRequest) (*Configuration, error) {
	configType := ConfigurationType(payload.Configuration.ConfigType)
	accountID := payload.Configuration.AccountID
	defaultConfig := false
	if payload.Configuration.DefaultConfiguration != nil {
		defaultConfig = *payload.Configuration.DefaultConfiguration
	}

	configData := payload.Configuration.ConfigurationData
	if configData == nil {
		configData = make(map[string]any)
	}

	cfg := &Configuration{
		AccountID:            accountID,
		DefaultConfiguration: defaultConfig,
		ConfigType:           string(configType),
		ConfigurationData:    configData,
	}
	if payload.Configuration.ID > 0 {
		cfg.ID = payload.Configuration.ID
	}
	return cfg, nil
}

func (c *Configuration) ToSMTPConfiguration() (*SMTPConfiguration, error) {
	jsonData, err := json.Marshal(c.ConfigurationData)
	if err != nil {
		return nil, err
	}
	var result SMTPConfiguration
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Configuration) ToWebAppConfiguration() (*WebAppConfiguration, error) {
	jsonData, err := json.Marshal(c.ConfigurationData)
	if err != nil {
		return nil, err
	}
	var result WebAppConfiguration
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Configuration) ToMap() (map[string]any, error) {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}
	return result, nil
}
