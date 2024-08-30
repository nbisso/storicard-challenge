package internal

import (
	"bytes"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/nbisso/storicard-challenge/domain"
)

type MigrationRepository interface {
	CreateMigration(ctx context.Context, m domain.Migration) (*domain.Migration, error)
	SaveMigrationFile(ctx context.Context, m domain.MigrationRequest) (string, error)
	GetMigrationFile(ctx context.Context, path string) ([]byte, error)
	UpdateMigration(ctx context.Context, m domain.Migration) error
	GetMigrationByFilename(ctx context.Context, filename string) (*domain.Migration, error)
	SaveTransaction(ctx context.Context, transaction domain.Transaction, file string) error
	GetUserBalance(ctx context.Context, tfilter domain.TransactionFilter) (domain.TransactionResult, error)
	GetFinishedMigrations(ctx context.Context) ([]domain.Migration, error)
	UpdateMigrationStatus(ctx context.Context, id int, status domain.MigrationStatus) error
}

type migrationRepository struct {
	db         sqlx.DB
	fileClient minio.Client
}

func NewMigrationRepository(db sqlx.DB, fileClient minio.Client) MigrationRepository {
	return &migrationRepository{
		db:         db,
		fileClient: fileClient,
	}
}

func (mr *migrationRepository) CreateMigration(ctx context.Context, m domain.Migration) (*domain.Migration, error) {
	_, err := mr.db.NamedExec(
		"INSERT INTO migration (csv_path, status,total_lines,processed_lines) VALUES (:csv_path, :status, :total_lines, :processed_lines)",
		m,
	)

	if err != nil {
		return nil, err
	}

	migrations := []domain.Migration{}
	err = mr.db.Select(&migrations, "SELECT * FROM migration WHERE csv_path=?", m.CsvPath)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &migrations[0], nil
}

func (mr *migrationRepository) SaveMigrationFile(ctx context.Context, m domain.MigrationRequest) (string, error) {

	uuidFileName := uuid.New().String()

	path := uuidFileName + ".csv"

	fileReader := bytes.NewReader(m.CsvFile)

	_, err := mr.fileClient.PutObject(ctx, "migrations", path, fileReader, int64(fileReader.Len()), minio.PutObjectOptions{})

	if err != nil {
		return "", err
	}

	return path, nil
}

func (mr *migrationRepository) GetMigrationFile(ctx context.Context, path string) ([]byte, error) {
	object, err := mr.fileClient.GetObject(ctx, "migrations", path, minio.GetObjectOptions{})

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(object)

	return buf.Bytes(), nil
}

func (mr *migrationRepository) UpdateMigration(ctx context.Context, m domain.Migration) error {
	_, err := mr.db.NamedExec(
		"UPDATE migration SET status=:status, total_lines=:total_lines, processed_lines=:processed_lines WHERE csv_path=:csv_path",
		m,
	)

	return err
}

func (mr *migrationRepository) GetMigrationByFilename(ctx context.Context, filename string) (*domain.Migration, error) {
	migrations := []domain.Migration{}
	err := mr.db.Select(&migrations, "SELECT * FROM migration WHERE csv_path=?", filename)

	if err != nil {
		return nil, err
	}

	return &migrations[0], nil
}

func (mr *migrationRepository) SaveTransaction(ctx context.Context, transaction domain.Transaction, file string) error {
	tx := mr.db.MustBegin()

	_, err := tx.NamedExec(
		"INSERT INTO transaction (user_id, amount, date_time) VALUES (:user_id, :amount, :datetime)",
		transaction,
	)

	if err != nil {
		tx.Rollback()
	}

	_, err = tx.Exec("UPDATE migration SET processed_lines=processed_lines+1 WHERE csv_path=?", file)

	if err != nil {
		tx.Rollback()
	}

	err = tx.Commit()

	return err

}

func (mr *migrationRepository) GetUserBalance(ctx context.Context, tfilter domain.TransactionFilter) (domain.TransactionResult, error) {
	transaction := []domain.TransactionResult{}

	userexists := false

	err := mr.db.Get(&userexists, "SELECT EXISTS(SELECT 1 FROM transaction WHERE user_id=?)", tfilter.UserID)

	if err != nil {
		return domain.TransactionResult{}, err
	}

	if !userexists {
		return domain.TransactionResult{}, domain.ErrUserNotFound
	}

	query := "SELECT SUM(amount) as balance, COUNT(CASE WHEN amount < 0 THEN amount ELSE NULL END) as total_debits, COUNT(CASE WHEN amount > 0 THEN amount ELSE NULL END) as total_credits FROM transaction WHERE user_id=?"

	args := []interface{}{tfilter.UserID}

	if tfilter.From != nil {
		query += " AND date_time >= ?"
		args = append(args, tfilter.From)
	}

	if tfilter.To != nil {
		query += " AND date_time <= ?"
		args = append(args, tfilter.To)
	}

	err = mr.db.Select(&transaction,
		query,
		args...)

	if err != nil {
		return domain.TransactionResult{}, err
	}

	return transaction[0], nil
}

func (mr *migrationRepository) GetFinishedMigrations(ctx context.Context) ([]domain.Migration, error) {
	migrations := []domain.Migration{}

	err := mr.db.Select(&migrations, "SELECT * FROM migration WHERE status != 'complete' AND total_lines = processed_lines")

	if err != nil {
		return nil, err
	}

	return migrations, nil
}

func (mr *migrationRepository) UpdateMigrationStatus(ctx context.Context, id int, status domain.MigrationStatus) error {
	_, err := mr.db.Exec("UPDATE migration SET status=? WHERE id=?", status, id)

	return err
}
