package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/nbisso/storicard-challenge/domain"
	queueclient "github.com/nbisso/storicard-challenge/infrastracture/queue_client"
)

type TransactionProcesssor interface {
	Start(ctx context.Context)
}

type transactionProcesssor struct {
	mu      MigrationUsecases
	qcTrans queueclient.QueueConsumerClient
}

func NewTransactionProcesssor(mu MigrationUsecases,
	qct queueclient.QueueConsumerClient,
) TransactionProcesssor {
	return &transactionProcesssor{
		mu:      mu,
		qcTrans: qct,
	}
}

func (f *transactionProcesssor) Start(ctx context.Context) {
	go func() {
		defer f.qcTrans.Close()
		for {
			msg, err := f.qcTrans.ReadMessage()
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

				transaction, err := domain.NewTransactionEventFromJson(string(msg.Value))

				if err != nil {
					log.Printf("Failed parse message: %s", err)
				}

				err = f.mu.SaveTransaction(ctx, *transaction)

				if err != nil {
					log.Printf("Failed to commit message: %s", err)
				}

				err = f.qcTrans.CommitMessage(msg)

				if err != nil {
					log.Printf("Failed to commit message: %s", err)
				}

				fmt.Printf("Message committed\n")

			} else {
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}()
}
