package apps

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/Scalingo/cli/api"
	"gopkg.in/errgo.v1"
)

var loadingRunes = "-\\|/"

func handleOperation(app string, res *http.Response) error {
	opURL, err := url.Parse(res.Header.Get("Location"))
	if err != nil {
		return errgo.Mask(err)
	}

	var op *api.Operation

	opID := filepath.Base(opURL.Path)
	done := make(chan struct{})
	errs := make(chan error)
	defer close(done)
	defer close(errs)

	go func() {
		for {
			op, err = api.OperationsShow(app, opID)
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
	for i := 0; ; i++ {
		select {
		case err := <-errs:
			return errgo.Mask(err)
		case <-done:
			if op.Status == "done" {
				fmt.Printf("\bDone in %.3f seconds\n", op.ElapsedDuration())
				return nil
			} else if op.Status == "error" {
				fmt.Printf("\bOperation '%s' failed, an error occured: %v\n", op.Type, op.Error)
				return nil
			}
		default:
		}

		r := loadingRunes[i%len(loadingRunes)]
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("\b%c", r)
	}

	return nil
}
