package registry

import (
	"context"

	"github.com/nbisso/storicard-challenge/infrastracture/conf"
	"github.com/nbisso/storicard-challenge/internal"
)

type Register interface {
	CleanUp()
}

type Registry struct {
	MigrationUsecases internal.MigrationUsecases
	Register          Register
}

type register struct {
}

func NewRegistry() *Registry {
	r := &register{}

	fp := r.NewFileProcessor()

	fp.Start(context.Background())

	tp := r.NewTransactionProcesssor()

	tp.Start(context.Background())

	status := r.NewStatusUpdater()

	status.Watch()

	ffp := r.NewFinishFileProcesssor()

	ffp.Start(context.Background())

	return &Registry{
		MigrationUsecases: r.NewMigrationUsecases(),
		Register:          r,
	}
}

func (r *register) NewMigrationUsecases() internal.MigrationUsecases {
	topic := conf.Instance.Kafka.FileTopic

	return internal.NewMigrationUsecases(
		r.NewMigrationRepository(),
		r.NewQueueSenderClient(topic),
	)
}

func (r *register) NewMigrationRepository() internal.MigrationRepository {
	return internal.NewMigrationRepository(*r.NewDatabase(), r.NewMinIOClient())
}

func (r *register) NewTransactionProcesssor() internal.TransactionProcesssor {
	return internal.NewTransactionProcesssor(
		r.NewMigrationUsecases(),
		r.NewQueueConsumerClient(conf.Instance.Kafka.EventTopic),
	)
}

func (r *register) NewFileProcessor() internal.FileProcessor {
	return internal.NewFileProcessor(r.NewMigrationUsecases(),
		r.NewQueueConsumerClient(conf.Instance.Kafka.FileTopic),
		r.NewQueueSenderClient(conf.Instance.Kafka.EventTopic))
}

func (r *register) NewFinishFileProcesssor() internal.FinishFileProcesssor {
	return internal.NewFinishFileProcesssor(r.NewMigrationUsecases(),
		r.NewQueueConsumerClient(conf.Instance.Kafka.FinishTopic),
	)
}

func (r *register) CleanUp() {
	r.CloseQueueClients()
}
