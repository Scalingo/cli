package apps

import (
	"bytes"
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/Scalingo/go-scalingo/v11"
)

func TestWaitOperationReturnsPollingErrorWithoutPanic(t *testing.T) {
	ctx := t.Context()

	expectedErr := stderrors.New("boom")
	callCount := 0

	waiter := &OperationWaiter{
		output: &bytes.Buffer{},
		prompt: defaultOperationWaiterPrompt,
		app:    "my-app",
		url:    "https://api.scalingo.test/operations/op-123",
		showOperation: func(_ context.Context, _ string, _ string) (*scalingo.Operation, error) {
			callCount++
			if callCount == 1 {
				return &scalingo.Operation{
					ID:     "op-123",
					Status: scalingo.OperationStatusRunning,
				}, nil
			}
			return nil, expectedErr
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

	if !strings.Contains(err.Error(), "get operation op-123") {
		t.Fatalf("expected operation id in error, got %v", err)
	}

	if !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Fatalf("expected wrapped polling error, got %v", err)
	}
}
