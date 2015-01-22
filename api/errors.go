package api

import (
	"fmt"
	"net/http"
	"strings"
)

type RequestFailedError struct {
	Code    int
	Message string
	Req     *http.Request
}

func NewRequestFailedError(code int, msg string, req *http.Request) *RequestFailedError {
	return &RequestFailedError{code, msg, req}
}

func (err *RequestFailedError) Error() string {
	return err.Message
}

func (err *RequestFailedError) String() string {
	return err.Message
}

func IsRequestFailedError(err error) bool {
	_, ok := err.(*RequestFailedError)
	return ok
}

type UnprocessableEntity struct {
	Errors map[string][]string `json:"errors"`
}

func (err *UnprocessableEntity) Error() string {
	errArray := make([]string, 0, len(err.Errors))
	for attr, attrErrs := range err.Errors {
		errArray = append(errArray, fmt.Sprintf("* %s → %s", attr, strings.Join(attrErrs, ", ")))
	}
	return strings.Join(errArray, "\n")
}
