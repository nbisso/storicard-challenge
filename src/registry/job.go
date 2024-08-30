package registry

import (
	"github.com/nbisso/storicard-challenge/infrastracture/conf"
	"github.com/nbisso/storicard-challenge/internal"
)

func (r *register) NewStatusUpdater() internal.StatusUpdater {

	topic := conf.Instance.Kafka.FinishTopic

	return internal.NewStatusUpdater(r.NewMigrationUsecases(), r.NewQueueSenderClient(topic))
}
