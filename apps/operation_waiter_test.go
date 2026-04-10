package apps

import (
	"bytes"
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/Scalingo/go-scalingo/v11"
)

type testOperationsService struct {
	show func(ctx context.Context, app, opID string) (*scalingo.Operation, error)
}

func (s testOperationsService) OperationsShow(ctx context.Context, app, opID string) (*scalingo.Operation, error) {
	return s.show(ctx, app, opID)
}

func TestWaitOperationReturnsPollingErrorWithoutPanic(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name               string
		expectedErr        error
		expectedSubstrings []string
	}{
		{
			name:        "return wrapped polling error",
			expectedErr: stderrors.New("boom"),
			expectedSubstrings: []string{
				"get operation op-123",
				"boom",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			callCount := 0
			waiter := &OperationWaiter{
				output: &bytes.Buffer{},
				prompt: defaultOperationWaiterPrompt,
				app:    "my-app",
				url:    "https://api.scalingo.test/operations/op-123",
				client: testOperationsService{
					show: func(_ context.Context, _ string, _ string) (*scalingo.Operation, error) {
						callCount++
						if callCount == 1 {
							return &scalingo.Operation{
								ID:     "op-123",
								Status: scalingo.OperationStatusRunning,
							}, nil
						}
						return nil, test.expectedErr
					},
				},
			}

			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("expected no panic, got %v", r)
				}
			}()

			_, err := waiter.WaitOperation(ctx)
			if err == nil {
				t.Fatal("expected polling error")
			}

			for _, expectedSubstring := range test.expectedSubstrings {
				if !strings.Contains(err.Error(), expectedSubstring) {
					t.Fatalf("expected %q in error, got %v", expectedSubstring, err)
				}
			}
		})
	}
}
