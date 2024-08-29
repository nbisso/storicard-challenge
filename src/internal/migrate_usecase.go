package internal

import (
	"context"

	"github.com/nbisso/storicard-challenge/domain"
	queueclient "github.com/nbisso/storicard-challenge/infrastracture/queue_client"
)

type MigrationUsecases interface {
	NewMigration(ctx context.Context, req domain.MigrationRequest) (domain.Migration, error)
	GetMigrationFile(ctx context.Context, filename string) (string, error)
	UpdateMigration(ctx context.Context, migration domain.Migration) error
	GetMigrationByFilename(ctx context.Context, filename string) (*domain.Migration, error)
	SaveTransaction(ctx context.Context, transaction domain.Transaction) error
	GetUserBalance(ctx context.Context, tfilter domain.TransactionFilter) (domain.TransactionResult, error)
}

type migrationUsecases struct {
	mr MigrationRepository
	qc queueclient.QueueSenderClient
}

func NewMigrationUsecases(mr MigrationRepository, qc queueclient.QueueSenderClient) MigrationUsecases {
	return &migrationUsecases{
		mr: mr,
		qc: qc,
	}
}

func (m *migrationUsecases) NewMigration(ctx context.Context, req domain.MigrationRequest) (domain.Migration, error) {

	newmigration := domain.Migration{
		Status: domain.Pending,
	}

	file, err := m.mr.SaveMigrationFile(ctx, req)

	if err != nil {
		newmigration.Status = domain.Failed
		newmigration.Summary = err.Error()
	}

	newmigration.CsvPath = file

	res, err := m.mr.CreateMigration(ctx, newmigration)

	if err != nil {
		res.Status = domain.Failed
		res.Summary = err.Error()
	}

	filevent := domain.NewFileEvent{
		FileName: res.CsvPath,
	}

	json, err := filevent.ToJson()

	if err != nil {
		res.Status = domain.Failed
		res.Summary = err.Error()
	}

	err = m.qc.SendMessage(json)

	if err != nil {
		res.Status = domain.Failed
		res.Summary = err.Error()
	}

	return *res, nil

}

func (m *migrationUsecases) GetMigrationFile(ctx context.Context, filename string) (string, error) {
	bytes, error := m.mr.GetMigrationFile(ctx, filename)

	if error != nil {
		return "", error
	}

	return string(bytes), nil
}

func (m *migrationUsecases) UpdateMigration(ctx context.Context, migration domain.Migration) error {
	return m.mr.UpdateMigration(ctx, migration)
}

func (m *migrationUsecases) GetMigrationByFilename(ctx context.Context, filename string) (*domain.Migration, error) {
	return m.mr.GetMigrationByFilename(ctx, filename)
}

func (m *migrationUsecases) SaveTransaction(ctx context.Context, transaction domain.Transaction) error {
	return m.mr.SaveTransaction(ctx, transaction)
}

func (m *migrationUsecases) GetUserBalance(ctx context.Context, tfilter domain.TransactionFilter) (domain.TransactionResult, error) {
	return m.mr.GetUserBalance(ctx, tfilter)
}
