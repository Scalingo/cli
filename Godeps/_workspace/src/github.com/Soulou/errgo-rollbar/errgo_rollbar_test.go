package errgorollbar

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/stvp/rollbar"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
)

func init() {
	rollbar.Environment = "test"
	rollbar.Token = os.Getenv("TOKEN")
}

func a() error {
	return errgo.New("error")
}

func b() error {
	return errgo.Mask(a())
}

func c() error {
	return errgo.Mask(b(), errgo.Any)
}

func note() error {
	return errgo.Notef(c(), "error of %s", "c")
}

func TestBuildStack(t *testing.T) {
	stack := BuildStack(note())
	if len(stack) != 4 {
		t.Errorf("len(stack) = %d != 4", len(stack))
	}
	fmt.Printf("%+v\n", stack)

	stack = BuildStack(c())
	if len(stack) != 3 {
		t.Errorf("len(stack) = %d != 3", len(stack))
	}

	stack = BuildStack(nil)
	if len(stack) != 0 {
		t.Errorf("len(stack) = %d != 0", len(stack))
	}

	stack = BuildStack(errors.New("error"))
	if len(stack) != 0 {
		t.Errorf("len(stack) = %d != 0", len(stack))
	}

	err := c()
	rollbar.ErrorWithStack(rollbar.ERR, err, BuildStack(err))
	rollbar.Wait()
}
