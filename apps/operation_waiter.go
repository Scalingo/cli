package apps

import (
	"context"
	"fmt"
	stdio "io"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

const (
	defaultOperationWaiterPrompt = "Status:  "
)

type OperationWaiter struct {
	output stdio.Writer
	prompt string
	app    string
	url    string
	client scalingo.OperationsService
}

func newOperationWaiterFromURL(ctx context.Context, app, url string) (*OperationWaiter, error) {
	return newOperationWaiter(ctx, os.Stderr, app, url)
}

func newOperationWaiter(
	ctx context.Context,
	output stdio.Writer,
	app,
	url string,
) (*OperationWaiter, error) {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "get Scalingo client")
	}

	return &OperationWaiter{
		output: output,
		app:    app,
		url:    url,
		prompt: defaultOperationWaiterPrompt,
		client: c,
	}, nil
}

func (w *OperationWaiter) SetPrompt(p string) {
	w.prompt = p
}

func (w *OperationWaiter) WaitOperation(ctx context.Context) (*scalingo.Operation, error) {
	opURL, err := url.Parse(w.url)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "parse url of operation")
	}

	var op *scalingo.Operation

	opID := filepath.Base(opURL.Path)
	done := make(chan struct{})
	errs := make(chan error)
	defer close(done)
	defer close(errs)

	op, err = w.client.OperationsShow(ctx, w.app, opID)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "get operation %v", opID)
	}

	go func() {
		for {
			nextOp, err := w.client.OperationsShow(ctx, w.app, opID)
			if err != nil {
				errs <- err
				break
			}
			op = nextOp

			if op.Status == "done" || op.Status == "error" {
				done <- struct{}{}
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Fprint(w.output, w.prompt)
	spinner := io.NewSpinner(os.Stderr)
	go spinner.Start()
	defer spinner.Stop()

	for {
		select {
		case err := <-errs:
			return op, errors.Wrapf(ctx, err, "get operation %v", opID)
		case <-done:
			switch op.Status {
			case "done":
				fmt.Printf("\bDone in %.3f seconds\n", op.ElapsedDuration())
				return op, nil
			case "error":
				fmt.Printf("\bOperation '%s' failed, an error occurred: %v\n", op.Type, op.Error)
				return op, errors.Newf(ctx, "operation %v failed", op.ID)
			}
		}
	}
}
