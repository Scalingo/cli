package dbng

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-scalingo/v11/scalingomock"
	"github.com/Scalingo/go-utils/retry"
)

func TestWaitForRunningDatabase(t *testing.T) {
	apiErr := stderrors.New("API unavailable")
	tests := map[string]struct {
		responses   []databaseShowResponse
		maxAttempts int
		wantErr     string
		wantIs      error
	}{
		"retries until the database is running": {
			responses: []databaseShowResponse{
				{database: databaseWithStatus(scalingo.DatabaseStatusCreating)},
				{database: databaseWithStatus(scalingo.DatabaseStatusMigrating)},
				{database: databaseWithStatus(scalingo.DatabaseStatusRunning)},
			},
			maxAttempts: 3,
		},
		"retries when the database is not visible yet": {
			responses: []databaseShowResponse{
				{err: scalingo.ErrDatabaseNotFound},
				{database: databaseWithStatus(scalingo.DatabaseStatusRunning)},
			},
			maxAttempts: 2,
		},
		"stops immediately when the database is stopped": {
			responses: []databaseShowResponse{
				{database: databaseWithStatus(scalingo.DatabaseStatusStopped)},
			},
			maxAttempts: 5,
			wantErr:     "database db-123 is stopped",
		},
		"stops immediately on an API error": {
			responses: []databaseShowResponse{
				{err: apiErr},
			},
			maxAttempts: 5,
			wantErr:     "get database db-123 status",
			wantIs:      apiErr,
		},
		"returns after the retry limit": {
			responses: []databaseShowResponse{
				{database: databaseWithStatus(scalingo.DatabaseStatusCreating)},
				{database: databaseWithStatus(scalingo.DatabaseStatusCreating)},
			},
			maxAttempts: 2,
			wantErr:     "database db-123 status is creating",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			previewClient := scalingomock.NewMockDatabasesPreviewService(ctrl)
			for _, response := range tt.responses {
				previewClient.EXPECT().DatabaseShow(gomock.Any(), "db-123").Return(response.database, response.err)
			}

			err := waitForRunningDatabase(t.Context(), previewClient, testRetryer(tt.maxAttempts), "db-123")

			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.wantErr)
			}
			if tt.wantIs != nil {
				require.ErrorIs(t, err, tt.wantIs)
			}
		})
	}

	// Ensure cancellation interrupts a status lookup that is still in progress.
	t.Run("cancels active lookup", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		previewClient := scalingomock.NewMockDatabasesPreviewService(ctrl)
		previewClient.EXPECT().DatabaseShow(gomock.Any(), "db-123").DoAndReturn(
			func(ctx context.Context, _ string) (scalingo.DatabaseNG, error) {
				deadline, ok := ctx.Deadline()
				require.True(t, ok)
				assert.WithinDuration(t, time.Now().Add(10*time.Millisecond), deadline, 10*time.Millisecond)

				<-ctx.Done()
				return scalingo.DatabaseNG{}, ctx.Err()
			},
		)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
		defer cancel()
		err := waitForRunningDatabase(ctx, previewClient, retry.New(retry.WithoutMaxAttempts()), "db-123")

		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.ErrorContains(t, err, "get database db-123 status")
	})
}

type databaseShowResponse struct {
	database scalingo.DatabaseNG
	err      error
}

func databaseWithStatus(status scalingo.DatabaseStatus) scalingo.DatabaseNG {
	return scalingo.DatabaseNG{Database: scalingo.Database{Status: status}}
}

func testRetryer(maxAttempts int) retry.Retry {
	return retry.New(
		retry.WithMaxAttempts(maxAttempts),
		retry.WithWaitDuration(time.Nanosecond),
	)
}
