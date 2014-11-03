package api

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type RequestFailedError struct {
	Message string
	Req     *http.Request
}

func NewRequestFailedError(msg string, req *http.Request) *RequestFailedError {
	return &RequestFailedError{msg, req}
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
	b := new(bytes.Buffer)
	for attr, attrErrs := range err.Errors {
		b.WriteString(fmt.Sprintf("%s → %s\n", attr, strings.Join(attrErrs, ", ")))
	}
	return b.String()
}
