package internal

import (
	"context"
	"fmt"
	"log"

	queueclient "github.com/nbisso/storicard-challenge/infrastracture/queue_client"
)

type FinishFileProcesssor interface {
	Start(ctx context.Context)
}

type finishFileProcesssor struct {
	mu      MigrationUsecases
	qcTrans queueclient.QueueConsumerClient
}

func NewFinishFileProcesssor(mu MigrationUsecases,
	qct queueclient.QueueConsumerClient,
) FinishFileProcesssor {
	return &finishFileProcesssor{
		mu:      mu,
		qcTrans: qct,
	}
}

func (f *finishFileProcesssor) Start(ctx context.Context) {
	go func() {
		defer f.qcTrans.Close()
		for {
			msg, err := f.qcTrans.ReadMessage()
			if err == nil {
				fmt.Printf("File finished to process on %s: %s\n", msg.TopicPartition, string(msg.Value))

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
