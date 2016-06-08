package errgorollbar

import (
	"github.com/stvp/rollbar"
	"gopkg.in/errgo.v1"
)

func BuildStack(err error) rollbar.Stack {
	stack := rollbar.Stack{}
	for err != nil {
		if errgoErr, ok := err.(*errgo.Err); !ok {
			break
		} else {
			frame := rollbar.Frame{
				Filename: errgoErr.File,
				Line:     errgoErr.Line,
			}
			stack = append([]rollbar.Frame{frame}, stack...)
			err = errgoErr.Underlying()
		}
	}
	return stack
}

// BuildStackWithSkip concatenates the stack given by the current execution flow
// and the stack determined by the errgo error
func BuildStackWithSkip(err error, skip int) rollbar.Stack {
	errStack := BuildStack(err)
	execStack := rollbar.BuildStack(skip + 2)
	return append(errStack, execStack...)
}
