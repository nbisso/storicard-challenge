package registry

import (
	"github.com/nbisso/storicard-challenge/infrastracture/conf"
	queueclient "github.com/nbisso/storicard-challenge/infrastracture/queue_client"
)

type Closable interface {
	Close()
}

var registeredQueueClient []Closable

func (r *register) NewQueueSenderClient(topic string) queueclient.QueueSenderClient {

	client := queueclient.NewQueueSenderClient(conf.Instance.Kafka.Host, topic)

	registeredQueueClient = append(registeredQueueClient, client)

	return client
}

func (r *register) NewQueueConsumerClient(topic string) queueclient.QueueConsumerClient {

	client := queueclient.NewQueueConsumerClient(conf.Instance.Kafka.Host, topic)

	registeredQueueClient = append(registeredQueueClient, client)

	return client
}

func (r *register) CloseQueueClients() {
	for _, client := range registeredQueueClient {
		client.Close()
	}
}
