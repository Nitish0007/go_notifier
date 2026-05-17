package rabbitmq

import (
	"os"
	"errors"

	"gopkg.in/yaml.v2"
	"github.com/Nitish0007/go_notifier/internal/common/mq"
)

const rabbitmqConfigPath = "configs/rabbitmq/rabbitmq.yml"

type envBlock struct {
	DynamicQueue bool     `yaml:"dynamic_queue"`
	Queues       []string `yaml:"queues"`
}

func InitializeQueues(mqClient mq.MQClient) ([]string, error) {
	env := os.Getenv("ENV")
	var configs map[string]envBlock
	yamlFile, err := os.ReadFile(rabbitmqConfigPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		return nil, err
	}

	block, ok := configs[env]
	if !ok {
		return nil, errors.New("invalid environment")
	}

	
	queues := block.Queues
	for _, queue := range queues {
		err := CreateQueue(mqClient, queue)
		if err != nil {
			return nil, err
		}
	}
	return block.Queues, nil
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