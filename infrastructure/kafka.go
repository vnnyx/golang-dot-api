package infrastructure

import (
	"context"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/vnnyx/golang-dot-api/model/message"
)

func createKafkaTopic() error {
	configuration := NewConfig(".env")
	a, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": configuration.BrokerHost,
	})
	if err != nil {
		return err
	}
	defer a.Close()

	_, err = a.CreateTopics(
		context.Background(),
		[]kafka.TopicSpecification{{
			Topic:             message.USER_OTP_TOPIC,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}},
	)
	if err != nil {
		return err
	}
	return nil
}

func NewKafkaProducer() (*kafka.Producer, error) {
	configuration := NewConfig(".env")
	hostName, _ := os.Hostname()
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": configuration.BrokerHost,
		"client.id":         hostName,
		"acks":              "all",
	})

	if err != nil {
		return nil, err
	}

	createKafkaTopic()

	return p, nil
}

func NewKafkaConsumer() (*kafka.Consumer, error) {
	configuration := NewConfig(".env")
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               configuration.BrokerHost,
		"group.id":                        "auth-service",
		"go.application.rebalance.enable": true,
	})
	if err != nil {
		return nil, err
	}

	createKafkaTopic()

	return consumer, nil
}
