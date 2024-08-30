package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/nbisso/storicard-challenge/domain"
	queueclient "github.com/nbisso/storicard-challenge/infrastracture/queue_client"
)

type StatusUpdater interface {
	Watch()
}

type statusUpdater struct {
	uc     MigrationUsecases
	sender queueclient.QueueSenderClient
}

func NewStatusUpdater(uc MigrationUsecases, qs queueclient.QueueSenderClient) StatusUpdater {
	return &statusUpdater{
		uc:     uc,
		sender: qs,
	}
}

func (s *statusUpdater) Watch() {
	go func() {
		ctx := context.Background()
		for {
			finishedMigrations, err := s.uc.GetFinishedMigrations(ctx)

			if err != nil {
				fmt.Println("Error getting finished migrations: ", err)
			}

			for _, migration := range finishedMigrations {

				err = s.uc.UpdateMigrationStatus(ctx, migration.Id, domain.Complete)

				if err != nil {
					fmt.Println("Error updating migration status: ", err)

					continue
				}

				err = s.sender.SendMessage(migration.CsvPath)

				if err != nil {
					fmt.Println("Error sending message: ", err)

					continue
				}

			}

			time.Sleep(5 * time.Second)
		}
	}()
}
