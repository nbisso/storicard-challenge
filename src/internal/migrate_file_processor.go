package internal

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/nbisso/storicard-challenge/domain"
	queueclient "github.com/nbisso/storicard-challenge/infrastracture/queue_client"
)

type FileProcessor interface {
	Start(ctx context.Context)
}

type fileProcessor struct {
	mu      MigrationUsecases
	qc      queueclient.QueueConsumerClient
	qcTrans queueclient.QueueSenderClient
}

func NewFileProcessor(mu MigrationUsecases, qc queueclient.QueueConsumerClient, qct queueclient.QueueSenderClient) FileProcessor {
	return &fileProcessor{
		mu:      mu,
		qc:      qc,
		qcTrans: qct,
	}
}

func (f *fileProcessor) Start(ctx context.Context) {
	go func() {
		defer f.qc.Close()
		for {
			msg, err := f.qc.ReadMessage()
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

				filereq, err := domain.NewNewFileEventFromJson(string(msg.Value))

				if err != nil {
					log.Printf("Failed to commit message: %s", err)
				}

				file, err := f.mu.GetMigrationFile(ctx, filereq.FileName)

				if err != nil {
					log.Printf("Failed to get file: %s", err)
				}

				transactions := []*domain.Transaction{}

				reader := gocsv.LazyCSVReader(strings.NewReader(string(file)))

				err = gocsv.UnmarshalCSV(reader, &transactions)

				if err != nil {
					log.Printf("Failed to unmarshal csv: %s", err)
				}

				m, err := f.mu.GetMigrationByFilename(ctx, filereq.FileName)

				if err != nil {
					log.Printf("Failed to get migration: %s", err)
				}

				m.Lines = len(transactions)
				m.Status = domain.Processing

				for _, migration := range transactions {
					j, err := migration.ToJson()

					if err != nil {
						log.Printf("Failed to convert to json: %s", err)
					}

					err = f.qcTrans.SendMessage(string(j))

					if err != nil {
						log.Printf("Failed to send message: %s", err)
					}
				}

				err = f.mu.UpdateMigration(ctx, *m)

				if err != nil {
					log.Printf("Failed to update migration: %s", err)
				}

				err = f.qc.CommitMessage(msg)

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
