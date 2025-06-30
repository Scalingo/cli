package apps

import (
	"context"
	"fmt"
	stdio "io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v8"
	errors "github.com/Scalingo/go-utils/errors/v2"
)

const (
	defaultOperationWaiterPrompt = "Status:  "
)

type OperationWaiter struct {
	output stdio.Writer
	prompt string
	app    string
	url    string
}

func NewOperationWaiterFromHTTPResponse(app string, res *http.Response) *OperationWaiter {
	operationURL := res.Header.Get("Location")
	return NewOperationWaiterFromURL(app, operationURL)
}

func NewOperationWaiterFromURL(app, url string) *OperationWaiter {
	return NewOperationWaiter(os.Stderr, app, url)
}

func NewOperationWaiter(output stdio.Writer, app, url string) *OperationWaiter {
	return &OperationWaiter{
		output: output,
		app:    app,
		url:    url,
		prompt: defaultOperationWaiterPrompt,
	}
}

func (w *OperationWaiter) SetPrompt(p string) {
	w.prompt = p
}

func (w *OperationWaiter) WaitOperation(ctx context.Context) (*scalingo.Operation, error) {
	opURL, err := url.Parse(w.url)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "parse url of operation")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "get Scalingo client")
	}

	var op *scalingo.Operation

	opID := filepath.Base(opURL.Path)
	done := make(chan struct{})
	errs := make(chan error)
	defer close(done)
	defer close(errs)

	op, err = c.OperationsShow(ctx, w.app, opID)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "get operation %v", opID)
	}

	go func() {
		for {
			op, err = c.OperationsShow(ctx, w.app, opID)
			if err != nil {
				errs <- err
				break
			}

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
			return op, errors.Wrapf(ctx, err, "get operation %v", op.ID)
		case <-done:
			if op.Status == "done" {
				fmt.Printf("\bDone in %.3f seconds\n", op.ElapsedDuration())
				return op, nil
			} else if op.Status == "error" {
				fmt.Printf("\bOperation '%s' failed, an error occurred: %v\n", op.Type, op.Error)
				return op, errors.Newf(ctx, "operation %v failed", op.ID)
			}
		}
	}
}
