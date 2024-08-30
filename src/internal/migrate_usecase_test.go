package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/nbisso/storicard-challenge/domain"
	mockqueue "github.com/nbisso/storicard-challenge/infrastracture/queue_client/mocks"
	"github.com/nbisso/storicard-challenge/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMigrationUsecases_NewMigration(t *testing.T) {
	var csvmock = `id,user_id,amount,datetime
59,4,42.47,2024-07-01T15:02:28.554Z
97,9,74.13,2024-07-01T13:46:19.916Z`
	type mockBehavior func(mr *mocks.MigrationRepository, qc *mockqueue.QueueSenderClient)

	cases := []struct {
		name          string
		req           domain.MigrationRequest
		mockBehavior  mockBehavior
		expectedError bool
		expectedRes   domain.Migration
	}{
		{
			name: "Success",
			req:  domain.MigrationRequest{},
			mockBehavior: func(mr *mocks.MigrationRepository, qc *mockqueue.QueueSenderClient) {
				mr.On("SaveMigrationFile", mock.Anything, mock.Anything).Return("path/to/file.csv", nil)
				mr.On("CreateMigration", mock.Anything, mock.Anything).Return(&domain.Migration{CsvPath: "path/to/file.csv", Status: domain.Pending}, nil)
				qc.On("SendMessage", mock.Anything).Return(nil)
			},
			expectedError: false,
			expectedRes: domain.Migration{
				CsvPath: "path/to/file.csv",
				Status:  domain.Pending,
			},
		},
		{
			name: "Failure when saving migration file",
			req: domain.MigrationRequest{
				CsvFile: []byte(csvmock),
			},
			mockBehavior: func(mr *mocks.MigrationRepository, qc *mockqueue.QueueSenderClient) {
				mr.On("SaveMigrationFile", mock.Anything, mock.Anything).Return("", errors.New("error saving file"))
			},
			expectedError: true,
			expectedRes: domain.Migration{
				Status:  domain.Failed,
				Summary: "error saving file",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			mr := new(mocks.MigrationRepository)
			qc := new(mockqueue.QueueSenderClient)

			tc.mockBehavior(mr, qc)

			uc := NewMigrationUsecases(mr, qc)

			res, err := uc.NewMigration(ctx, tc.req)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedRes.Status, res.Status)
				assert.Equal(t, tc.expectedRes.Summary, res.Summary)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}

			mr.AssertExpectations(t)
			qc.AssertExpectations(t)
		})
	}
}
