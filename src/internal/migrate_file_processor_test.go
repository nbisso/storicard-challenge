package internal

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/nbisso/storicard-challenge/domain"
	mockqueue "github.com/nbisso/storicard-challenge/infrastracture/queue_client/mocks"
	"github.com/nbisso/storicard-challenge/internal/mocks"
	"github.com/stretchr/testify/mock"
)

func TestFileProcessor_Start(t *testing.T) {
	type mockBehavior func(mr *mocks.MigrationUsecases, qc *mockqueue.QueueConsumerClient, qct *mockqueue.QueueSenderClient) *mockqueue.QueueSenderClientTransactioner
	var csvmock = `id,user_id,amount,datetime
59,4,42.47,2024-07-01T15:02:28.554Z
97,9,74.13,2024-07-01T13:46:19.916Z`
	cases := []struct {
		name          string
		mockBehavior  mockBehavior
		expectedError bool
	}{
		{
			name: "Success",
			mockBehavior: func(mu *mocks.MigrationUsecases, qc *mockqueue.QueueConsumerClient, qct *mockqueue.QueueSenderClient) *mockqueue.QueueSenderClientTransactioner {
				transaction := new(mockqueue.QueueSenderClientTransactioner)
				testMessage := domain.NewFileEvent{
					FileName: "test.csv",
				}
				messageJSON, _ := testMessage.ToJson()
				qc.On("ReadMessage").Return(&kafka.Message{Value: []byte(messageJSON)}, nil)
				qc.On("CommitMessage", mock.Anything).Return(nil)
				qct.On("InitTransaction", mock.Anything).Return(transaction, nil)
				transaction.On("BeginTransaction").Return(nil)
				transaction.On("SendMessageWithHeaders", mock.Anything, mock.Anything).Return(nil)
				transaction.On("CommitTransaction", mock.Anything).Return(nil)
				mu.On("GetMigrationFile", mock.Anything, "test.csv").Return(csvmock, nil)
				mu.On("GetMigrationByFilename", mock.Anything, "test.csv").Return(&domain.Migration{}, nil)
				mu.On("UpdateMigration", mock.Anything, mock.Anything).Return(nil)

				return transaction
			},
			expectedError: false,
		},
		{
			name: "Failure on UpdateMigration",
			mockBehavior: func(mu *mocks.MigrationUsecases, qc *mockqueue.QueueConsumerClient, qct *mockqueue.QueueSenderClient) *mockqueue.QueueSenderClientTransactioner {
				transaction := new(mockqueue.QueueSenderClientTransactioner)
				testMessage := domain.NewFileEvent{
					FileName: "test.csv",
				}
				messageJSON, _ := testMessage.ToJson()
				qc.On("ReadMessage").Return(&kafka.Message{Value: []byte(messageJSON)}, nil)
				qc.On("CommitMessage", mock.Anything).Return(nil)
				qct.On("InitTransaction", mock.Anything).Return(transaction, nil)
				transaction.On("BeginTransaction").Return(nil)
				transaction.On("SendMessageWithHeaders", mock.Anything, mock.Anything).Return(nil)
				transaction.On("CommitTransaction", mock.Anything).Return(nil)
				mu.On("GetMigrationFile", mock.Anything, "test.csv").Return(csvmock, nil)
				mu.On("GetMigrationByFilename", mock.Anything, "test.csv").Return(&domain.Migration{}, nil)
				mu.On("UpdateMigration", mock.Anything, mock.Anything).Return(errors.New("error updating migration"))
				transaction.On("RollbackTransaction", mock.Anything).Return(nil)

				return transaction
			},
			expectedError: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			mu := new(mocks.MigrationUsecases)
			qc := new(mockqueue.QueueConsumerClient)
			qs := new(mockqueue.QueueSenderClient)

			sender := tc.mockBehavior(mu, qc, qs)

			fp := NewFileProcessor(mu, qc, qs)

			go fp.Start(ctx)

			time.Sleep(2 * time.Second)

			if tc.expectedError {
				sender.AssertCalled(t, "RollbackTransaction", mock.Anything)
			} else {
				sender.AssertCalled(t, "CommitTransaction", mock.Anything)
			}

		})
	}

}
