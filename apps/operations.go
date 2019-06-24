package apps

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func handleOperation(app string, res *http.Response) error {
	opURL, err := url.Parse(res.Header.Get("Location"))
	if err != nil {
		return errgo.Mask(err)
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	var op *scalingo.Operation

	opID := filepath.Base(opURL.Path)
	done := make(chan struct{})
	errs := make(chan error)
	defer close(done)
	defer close(errs)

	go func() {
		for {
			op, err = c.OperationsShow(app, opID)
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

	fmt.Print("Status:  ")
	spinner := io.NewSpinner(os.Stderr)
	go spinner.Start()
	defer spinner.Stop()

	for {
		select {
		case err := <-errs:
			return errgo.Mask(err)
		case <-done:
			if op.Status == "done" {
				fmt.Printf("\bDone in %.3f seconds\n", op.ElapsedDuration())
				return nil
			} else if op.Status == "error" {
				fmt.Printf("\bOperation '%s' failed, an error occurred: %v\n", op.Type, op.Error)
				return nil
			}
		}
	}
}
