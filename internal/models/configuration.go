package models

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
	"github.com/Nitish0007/go_notifier/internal/validators"
)

type ConfigurationType string

const (
	SMTPConfig   ConfigurationType = "smtp"
	WebAppConfig ConfigurationType = "in_app" // this is for web app notifications
)

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

type Configuration struct {
	ID                   int            `json:"id" gorm:"primaryKey" validate:"-"`
	AccountID            int            `json:"account_id" gorm:"not null;index" validate:"required,gt=0"`
	DefaultConfiguration bool           `json:"default_configuration" gorm:"default:false" validate:"-"`
	ConfigurationData    map[string]any `json:"configuration_data" gorm:"type:jsonb;default:'{}'::jsonb;serializer:json" validate:"required"`
	ConfigType           string         `json:"config_type" gorm:"not null" validate:"required,oneof=smtp in_app"`
	CreatedAt            time.Time      `json:"created_at" gorm:"autoCreateTime" validate:"-"`
	UpdatedAt            time.Time      `json:"updated_at" gorm:"autoUpdateTime" validate:"-"`
}

func (c *Configuration) BeforeSave(tx *gorm.DB) error {
	log.Printf(">>>>>>>>>>>>>> before save hook\n")
	log.Printf("config: %v\n", c)
	if c.ConfigType == "" {
		return errors.New("config type is required")
	}

	cType := ConfigurationType(c.ConfigType)

	switch cType {
	case SMTPConfig:
		validator := validators.NewModelValidator[SMTPConfiguration]()
		smtpConfig, err := c.ToSMTPConfiguration()
		if err != nil {
			return err
		}
		if err := validator.ValidateStruct(smtpConfig); err != nil {
			return err
		}
	case WebAppConfig:
		validator := validators.NewModelValidator[WebAppConfiguration]()
		webAppConfig, err := c.ToWebAppConfiguration()
		if err != nil {
			return err
		}
		if err := validator.ValidateStruct(webAppConfig); err != nil {
			return err
		}
	}
	return nil
}

func (c *Configuration) BeforeCreate(tx *gorm.DB) error {
	if c.ConfigType == "" {
		return errors.New("config type is required")
	}

	if c.ConfigType == string(SMTPConfig) {
		smtpConfig, err := c.ToSMTPConfiguration()
		if err != nil {
			return err
		}
		c.ConfigurationData = map[string]any{
			"host":     smtpConfig.Host,
			"port":     smtpConfig.Port,
			"username": smtpConfig.Username,
			"password": smtpConfig.Password,
			"from":     smtpConfig.From,
		}
	} else if c.ConfigType == string(WebAppConfig) {
		webAppConfig, err := c.ToWebAppConfiguration()
		if err != nil {
			return err
		}
		c.ConfigurationData = map[string]any{
			"web_app_url":    webAppConfig.WebAppURL,
			"web_app_secret": webAppConfig.WebAppSecret,
			"web_app_token":  webAppConfig.WebAppToken,
		}
	}
	return nil
}

func NewConfiguration(configData map[string]any) (*Configuration, error) {
	configType, exists := configData["config_type"].(string)
	if !exists {
		return nil, errors.New("config type is required")
	}

	accountID, ok := configData["account_id"].(int)
	if !ok {
		// Try float64 (JSON numbers are often unmarshaled as float64)
		if accountIDFloat, ok := configData["account_id"].(float64); ok {
			accountID = int(accountIDFloat)
		} else {
			return nil, errors.New("account_id is required and must be a number")
		}
	}

	defaultConfig := false
	if dc, ok := configData["default_configuration"].(bool); ok {
		defaultConfig = dc
	}

	// Extract configuration_data if present, otherwise use the whole configData
	configDataMap, ok := configData["configuration_data"].(map[string]any)
	if !ok {
		// If configuration_data is not present, try to extract it from the root
		configDataMap = make(map[string]any)
		for k, v := range configData {
			if k != "account_id" && k != "config_type" && k != "default_configuration" && k != "id" {
				configDataMap[k] = v
			}
		}
	}

	config := &Configuration{
		AccountID:            accountID,
		DefaultConfiguration: defaultConfig,
		ConfigType:           configType,
		ConfigurationData:    configDataMap,
	}

	return config, nil
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

func (c *Configuration) ToConfiguration() (*Configuration, error) {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	var result Configuration
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func ValidateConfigs(configData map[string]any) error {
	if len(configData) == 0 {
		return errors.New("config data is required")
	}

	jsonData, err := json.Marshal(configData)
	if err != nil {
		return err
	}
	var result Configuration
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return err
	}
	if result.ConfigType == "" {
		return errors.New("config type is required")
	}
	return nil
}
