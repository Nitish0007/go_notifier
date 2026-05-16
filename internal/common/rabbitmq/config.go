package rabbitmq

import (
	"os"
	"errors"

	"gopkg.in/yaml.v2"
	"github.com/Nitish0007/go_notifier/internal/common/mq"
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

func InitializeQueues(mqClient mq.MQClient) ([]string, error) {
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
		queues := envConfig[env].Development.Queues
		for _, queue := range queues {
			err := CreateQueue(mqClient, queue)
			if err != nil {
				return nil, err
			}
		}
		return queues, nil
	}

	return nil, errors.New("invalid environment")
}

func CreateQueue(rbmqClient mq.MQClient, queue_name string) error {
	_, err := rbmqClient.(*RabbitMQClient).channel.QueueDeclare(
		queue_name, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	return err
}