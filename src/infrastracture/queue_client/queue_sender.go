package queueclient

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

type QueueSenderClientTransactioner interface {
	QueueSenderClient
	BeginTransaction() error
	CommitTransaction(ctx context.Context) error
	RollbackTransaction(ctx context.Context) error
}

type QueueSenderClient interface {
	SendMessage(message string) error
	SendMessageWithHeaders(message string, headers map[string]string) error
	Flush()
	Close()
	InitTransaction(ctx context.Context) (QueueSenderClientTransactioner, error)
}

type kafkaSenderTransactionerClient struct {
	transactionID string
	kafkaSenderClient
}

type kafkaSenderClient struct {
	Producer *kafka.Producer
	Topic    string
	config   kafka.ConfigMap
}

func (k *kafkaSenderTransactionerClient) BeginTransaction() error {

	return k.Producer.BeginTransaction()

}

func (k *kafkaSenderTransactionerClient) CommitTransaction(ctx context.Context) error {
	return k.Producer.CommitTransaction(ctx)

}

func (k *kafkaSenderTransactionerClient) RollbackTransaction(ctx context.Context) error {
	return k.Producer.AbortTransaction(ctx)
}

func (k *kafkaSenderClient) InitTransaction(ctx context.Context) (QueueSenderClientTransactioner, error) {
	transactionid := uuid.New().String()

	client := NewQueueTransactionalSenderClient(ctx, k.config["bootstrap.servers"].(string), k.Topic, transactionid)

	return client, nil
}

func (k *kafkaSenderClient) Flush() {
	k.Producer.Flush(15 * 1000)
}

func (k *kafkaSenderClient) internalSendMessageWithRetry(message string, headers map[string]string) error {
	retry := 0
	headerkafka := make([]kafka.Header, 0)

	var err error

	for k, v := range headers {
		headerkafka = append(headerkafka, kafka.Header{
			Key:   k,
			Value: []byte(v),
		})
	}

	for retry < 3 {

		err = k.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &k.Topic,
				Partition: kafka.PartitionAny,
			},
			Value:   []byte(message),
			Headers: headerkafka,
		},

			nil,
		)

		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrQueueFull {
				k.Producer.Flush(15 * 1000)
			}

		} else {
			break
		}

		retry++
	}

	return err
}

func (k *kafkaSenderClient) SendMessage(message string) error {

	empty := make(map[string]string)

	return k.internalSendMessageWithRetry(message, empty)
}

func NewQueueSenderClient(serverhost string, topic string) QueueSenderClient {

	config := kafka.ConfigMap{
		"bootstrap.servers":            serverhost,
		"queue.buffering.max.messages": 1000000,
		"queue.buffering.max.kbytes":   104857,
		"queue.buffering.max.ms":       1000,
		"batch.num.messages":           1000,
		"linger.ms":                    100,
		"compression.codec":            "snappy",
	}

	producer, err := kafka.NewProducer(&config)

	if err != nil {
		panic(err)
	}

	return &kafkaSenderClient{
		Producer: producer,
		Topic:    topic,
		config:   config,
	}
}

func NewQueueTransactionalSenderClient(ctx context.Context, serverhost string, topic string, transactionid string) QueueSenderClientTransactioner {

	config := kafka.ConfigMap{
		"bootstrap.servers":            serverhost,
		"queue.buffering.max.messages": 1000000,
		"queue.buffering.max.kbytes":   104857,
		"queue.buffering.max.ms":       1000,
		"batch.num.messages":           1000,
		"linger.ms":                    100,
		"compression.codec":            "snappy",
		"transactional.id":             transactionid,
	}

	producer, err := kafka.NewProducer(&config)

	if err != nil {
		panic(err)
	}

	err = producer.InitTransactions(ctx)

	if err != nil {
		panic(err)
	}

	return &kafkaSenderTransactionerClient{
		transactionID: transactionid,
		kafkaSenderClient: kafkaSenderClient{
			Producer: producer,
			Topic:    topic,
			config:   config,
		},
	}
}

func (k *kafkaSenderClient) Close() {
	k.Flush()
	k.Producer.Close()
}

func (k *kafkaSenderClient) SendMessageWithHeaders(message string, headers map[string]string) error {
	return k.internalSendMessageWithRetry(message, headers)
}
