package queueclient

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type QueueConsumerClient interface {
	CommitMessage(m *kafka.Message) error
	ReadMessage() (*kafka.Message, error)
	Close()
}

type consumerClient struct {
	Consumer *kafka.Consumer
	Topic    string
}

// CommitMessage implements QueueClient.
func (k *consumerClient) CommitMessage(m *kafka.Message) error {
	_, err := k.Consumer.CommitMessage(m)

	return err
}

// ReadMessage implements QueueClient.
func (k *consumerClient) ReadMessage() (*kafka.Message, error) {

	return k.Consumer.ReadMessage(-1)
}

func NewQueueConsumerClient(serverhost string, topic string) QueueConsumerClient {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     serverhost,
		"group.id":              "consumer-" + topic,
		"auto.offset.reset":     "earliest",
		"enable.auto.commit":    false,
		"session.timeout.ms":    6000, // Reduced from 10000 to 6000
		"heartbeat.interval.ms": 2000, // Reduced from 3000 to 2000
	})

	if err != nil {
		panic(err)
	}

	err = c.Subscribe(topic, nil)

	if err != nil {
		panic(err)
	}

	return &consumerClient{
		Topic:    topic,
		Consumer: c,
	}
}

func (k *consumerClient) Close() {
	k.Consumer.Unsubscribe()
}
