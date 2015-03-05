package api

import (
	"fmt"
	"net/http"
	"strings"
)

type PaymentRequiredError struct {
	ErrMessage string `json:"error"`
	URL        string `json:"url"`
}

func (err *PaymentRequiredError) Error() string {
	return fmt.Sprintf("%v\n→ %v", err.ErrMessage, err.URL)
}

type NotFoundError struct {
	Resource string `json:"resource"`
	Err      string `json:"error"`
}

func (err *NotFoundError) Error() string {
	if err.Resource == "app" {
		return fmt.Sprintf("The application has not been found, have you done a typo?")
	} else if err.Resource == "container_type" {
		return fmt.Sprintf("This type of container has not been found, please ensure it is present in your Procfile\n→ http://doc.scalingo.com/internals/procfile")
	} else {
		return fmt.Sprintf("The %s has not been found", err.Resource)
	}
}

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
