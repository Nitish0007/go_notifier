package rabbitmq

import (
	"os"
	"errors"

	"gopkg.in/yaml.v2"
)

const rabbitmqConfigPath = "configs/rabbitmq/rabbitmq.yml"

type envConfig struct {
	Development struct {
		DynamicQueue bool `yaml:"dynamic_queue"`
		Queues []string `yaml:"queues"`
	}
	Test struct {
		DynamicQueue bool `yaml:"dynamic_queue"`
		Queues []string `yaml:"queues"`
	}
	Production struct {
		DynamicQueue bool `yaml:"dynamic_queue"`
		Queues []string `yaml:"queues"`
	}
}

func InitializeQueues() ([]string, error) {
	env := os.Getenv("ENV")
	envConfig := make(map[string]envConfig)
	yamlFile, err := os.ReadFile(rabbitmqConfigPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &envConfig)
	if err != nil {
		return nil, err
	}

	switch env {
	case "development", "test", "production":
		if envConfig[env].Development.DynamicQueue == true {
			return []string{}, nil
		}
		return envConfig[env].Development.Queues, nil
	}

	return nil, errors.New("invalid environment")
}
