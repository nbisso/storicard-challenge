package domain

import "strings"

type MigrationStatus string

const (
	Pending    MigrationStatus = "pending"
	Complete   MigrationStatus = "complete"
	Failed     MigrationStatus = "failed"
	Processing MigrationStatus = "processing"
)

type Migration struct {
	Id             int             `db:"id" json:"id"`
	CsvPath        string          `db:"csv_path" json:"csv_path"`
	Lines          int             `db:"total_lines" json:"total_lines"`
	ProcessedLines int             `db:"processed_lines" json:"processed_lines"`
	Status         MigrationStatus `db:"status" json:"status"`
	Summary        string          `db:"summary" json:"summary"`
	CreatedAt      string          `db:"created_at" json:"created_at"`
	UpdatedAt      string          `db:"updated_at" json:"updated_at"`
}

type MigrationRequest struct {
	CsvFile []byte
	lines   *int
}

func (m *MigrationRequest) GetLines() int {

	if m.lines != nil {
		return *m.lines
	}

	lines := 0

	strginFile := string(m.CsvFile)

	lines = len(strings.Split(strginFile, "\n"))

	m.lines = &lines

	return *m.lines
}
