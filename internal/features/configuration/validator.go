package configuration

import (
	"fmt"

	"github.com/Nitish0007/go_notifier/internal/shared/validators"
)

// ValidateConfigurationDataByType validates flat configuration_data by config_type.
// Returns the validated flat map for storage. Adding a new type = add case here + struct + oneof in DTOs.
func ValidateConfigurationDataByType(configType string, data map[string]any) (map[string]any, error) {
	if data == nil {
		return nil, fmt.Errorf("configuration_data is required")
	}
	switch configType {
	case string(SMTPConfig):
		v := validators.NewModelValidator[SMTPConfiguration]()
		cfg, err := v.ValidateFromMap(data)
		if err != nil {
			return nil, err
		}
		flat, err := v.ToMap(cfg)
		if err != nil {
			return nil, err
		}
		return flat, nil
	case string(WebAppConfig):
		v := validators.NewModelValidator[WebAppConfiguration]()
		cfg, err := v.ValidateFromMap(data)
		if err != nil {
			return nil, err
		}
		flat, err := v.ToMap(cfg)
		if err != nil {
			return nil, err
		}
		return flat, nil
	default:
		return nil, fmt.Errorf("unsupported config_type: %s", configType)
	}
}
