package apps

import (
	"bytes"
	stderrors "errors"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-scalingo/v11/scalingomock"
)

func TestWaitOperationReturnsPollingErrorWithoutPanic(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	expectedErr := stderrors.New("boom")
	expectedSubstrings := []string{
		"get operation op-123",
		"boom",
	}
	operationsService := scalingomock.NewMockOperationsService(ctrl)
	operationsService.EXPECT().OperationsShow(ctx, "my-app", "op-123").Return(&scalingo.Operation{
		ID:     "op-123",
		Status: scalingo.OperationStatusRunning,
	}, nil)
	operationsService.EXPECT().OperationsShow(ctx, "my-app", "op-123").Return(nil, expectedErr)

	waiter := &OperationWaiter{
		output:            &bytes.Buffer{},
		prompt:            defaultOperationWaiterPrompt,
		app:               "my-app",
		url:               "https://api.scalingo.test/operations/op-123",
		operationsService: operationsService,
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

	for _, expectedSubstring := range expectedSubstrings {
		if !strings.Contains(err.Error(), expectedSubstring) {
			t.Fatalf("expected %q in error, got %v", expectedSubstring, err)
		}
	}
}
