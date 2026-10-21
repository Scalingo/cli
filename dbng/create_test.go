package dbng

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/retry"
)

func TestWaitForRunningDatabaseWithRetry(t *testing.T) {
	t.Run("retries until the database is running", func(t *testing.T) {
		calls := 0
		statuses := []scalingo.DatabaseStatus{
			scalingo.DatabaseStatusCreating,
			scalingo.DatabaseStatusMigrating,
			scalingo.DatabaseStatusRunning,
		}
		show := func(_ context.Context, appID string) (scalingo.DatabaseNG, error) {
			assert.Equal(t, "db-123", appID)
			status := statuses[calls]
			calls++
			return databaseWithStatus(status), nil
		}

		err := waitForRunningDatabaseWithRetry(t.Context(), show, testRetryer(3), "db-123")

		require.NoError(t, err)
		assert.Equal(t, 3, calls)
	})

	t.Run("retries a database that is not visible yet", func(t *testing.T) {
		calls := 0
		show := func(_ context.Context, _ string) (scalingo.DatabaseNG, error) {
			calls++
			if calls == 1 {
				return scalingo.DatabaseNG{}, scalingo.ErrDatabaseNotFound
			}
			return databaseWithStatus(scalingo.DatabaseStatusRunning), nil
		}

		err := waitForRunningDatabaseWithRetry(t.Context(), show, testRetryer(2), "db-123")

		require.NoError(t, err)
		assert.Equal(t, 2, calls)
	})

	t.Run("stops immediately when the database is stopped", func(t *testing.T) {
		calls := 0
		show := func(_ context.Context, _ string) (scalingo.DatabaseNG, error) {
			calls++
			return databaseWithStatus(scalingo.DatabaseStatusStopped), nil
		}

		err := waitForRunningDatabaseWithRetry(t.Context(), show, testRetryer(5), "db-123")

		require.ErrorContains(t, err, "database db-123 is stopped")
		assert.Equal(t, 1, calls)
	})

	t.Run("stops immediately on an API error", func(t *testing.T) {
		apiErr := stderrors.New("API unavailable")
		calls := 0
		show := func(_ context.Context, _ string) (scalingo.DatabaseNG, error) {
			calls++
			return scalingo.DatabaseNG{}, apiErr
		}

		err := waitForRunningDatabaseWithRetry(t.Context(), show, testRetryer(5), "db-123")

		require.ErrorIs(t, err, apiErr)
		require.ErrorContains(t, err, "get database db-123 status")
		assert.Equal(t, 1, calls)
	})

	t.Run("returns after the retry limit", func(t *testing.T) {
		calls := 0
		show := func(_ context.Context, _ string) (scalingo.DatabaseNG, error) {
			calls++
			return databaseWithStatus(scalingo.DatabaseStatusCreating), nil
		}

		err := waitForRunningDatabaseWithRetry(t.Context(), show, testRetryer(2), "db-123")

		require.ErrorContains(t, err, "database db-123 status is creating")
		assert.Equal(t, 2, calls)
	})
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
