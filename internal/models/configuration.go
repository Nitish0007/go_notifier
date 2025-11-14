package models

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ConfigurationType string

const (
	SMTPConfig 		ConfigurationType = "smtp"
	WebAppConfig 	ConfigurationType = "in_app" // this is for web app notifications
)

type SMTPConfiguration struct {
	Host string `json:"host"`
	Port int `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From string `json:"from"`
}

type WebAppConfiguration struct {
	WebAppURL    string `json:"web_app_url"`
	WebAppSecret string `json:"web_app_secret"`
	WebAppToken  string `json:"web_app_token"`
	WebAppPort   int `json:"web_app_port"`
	WebAppHost   string `json:"web_app_host"`
}


type Configuration struct {
	ID 											int 						`json:"id" gorm:"primaryKey"`
	AccountID 							int 						`json:"account_id" gorm:"not null;index"`
	DefaultConfiguration 		bool 						`json:"default_configuration" gorm:"default:false"`
	ConfigurationData 			map[string]any 	`json:"configuration_data" gorm:"type:jsonb;default:'{}'::jsonb"`
	ConfigType 							string 					`json:"config_type" gorm:"not null"`
	CreatedAt 							time.Time 			`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt 							time.Time 			`json:"updated_at" gorm:"autoUpdateTime"`
}

func (c *Configuration) BeforeCreate(tx *gorm.DB) error {
	if c.ConfigType == "" {
		return errors.New("config type is required")
	}

	jsonData, err := json.Marshal(c.ConfigurationData)
	if err != nil {
		return err
	}
	if c.ConfigType == string(SMTPConfig) {
		smtpConfig := &SMTPConfiguration{}
		if err := json.Unmarshal(jsonData, smtpConfig); err != nil {
			return err
		}
		c.ConfigurationData = map[string]any{
			"host": smtpConfig.Host,
			"port": smtpConfig.Port,
			"username": smtpConfig.Username,
			"password": smtpConfig.Password,
			"from": smtpConfig.From,
		}
	} else if c.ConfigType == string(WebAppConfig) {
		webAppConfig := &WebAppConfiguration{}
		if err := json.Unmarshal(jsonData, webAppConfig); err != nil {
			return err
		}
		c.ConfigurationData = map[string]any{
			"web_app_url": webAppConfig.WebAppURL,
			"web_app_secret": webAppConfig.WebAppSecret,
			"web_app_token": webAppConfig.WebAppToken,
		}
	}
	return nil
}

func NewConfiguration(configData map[string]any) (*Configuration, error) {
	
	configType, exists := configData["config_type"].(string)
	if !exists {
		return nil, errors.New("config type is required")
	}
	config := &Configuration{
		AccountID: configData["account_id"].(int),
		DefaultConfiguration: configData["default_configuration"].(bool),
		ConfigType: configType,
		ConfigurationData: configData["configuration_data"].(map[string]any),
	}
	if err := config.BeforeCreate(nil); err != nil {
		return nil, err
	}
	return config, nil
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