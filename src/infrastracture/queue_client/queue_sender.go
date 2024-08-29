package queueclient

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type QueueSenderClient interface {
	SendMessage(message string) error
	SendMessageWithHeaders(message string, headers map[string]string) error
	Flush()
	Close()
}

type kafkaSenderClient struct {
	Producer *kafka.Producer
	Consumer *kafka.Consumer
	Topic    string
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

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":            serverhost,
		"queue.buffering.max.messages": 1000000,
		"queue.buffering.max.kbytes":   104857,
		"queue.buffering.max.ms":       1000,
		"batch.num.messages":           1000,
		"linger.ms":                    100,
		"compression.codec":            "snappy",
	})

	if err != nil {
		panic(err)
	}

	return &kafkaSenderClient{
		Producer: producer,
		Topic:    topic,
	}
}

func (k *kafkaSenderClient) Close() {
	k.Flush()
	k.Producer.Close()
	k.Consumer.Unsubscribe()
}

func (k *kafkaSenderClient) SendMessageWithHeaders(message string, headers map[string]string) error {
	return k.internalSendMessageWithRetry(message, headers)
}
