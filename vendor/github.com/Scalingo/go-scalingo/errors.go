package scalingo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Scalingo/go-scalingo/debug"
	"gopkg.in/errgo.v1"
)

type (
	BadRequestError struct {
		ErrMessage string `json:"error"`
	}

	PaymentRequiredError struct {
		Name       string `json:"name"`
		ErrMessage string `json:"error"`
		URL        string `json:"url"`
	}

	NotFoundError struct {
		Resource string `json:"resource"`
		Err      string `json:"error"`
	}

	UnprocessableEntity struct {
		Errors map[string][]string `json:"errors"`
	}

	RequestFailedError struct {
		Code     int
		APIError error
		Req      *APIRequest
	}
)

func (err BadRequestError) Error() string {
	return fmt.Sprintf("400 Bad Request → %v", err.ErrMessage)
}

func (err PaymentRequiredError) Error() string {
	return fmt.Sprintf("%v\n→ %v", err.ErrMessage, err.URL)
}

func (err NotFoundError) Error() string {
	if err.Resource == "app" {
		return fmt.Sprintf("The application has not been found, have you done a typo?")
	} else if err.Resource == "container_type" {
		return fmt.Sprintf("This type of container has not been found, please ensure it is present in your Procfile\n→ http://doc.scalingo.com/internals/procfile")
	} else {
		return fmt.Sprintf("The %s has not been found", err.Resource)
	}
}

func (err UnprocessableEntity) Error() string {
	errArray := make([]string, 0, len(err.Errors))
	for attr, attrErrs := range err.Errors {
		errArray = append(errArray, fmt.Sprintf("* %s → %s", attr, strings.Join(attrErrs, ", ")))
	}
	return strings.Join(errArray, "\n")
}

func NewRequestFailedError(res *http.Response, req *APIRequest) error {
	debug.Println("APIRequest Error:", res.StatusCode, req.Method, req.Endpoint)
	defer res.Body.Close()
	switch res.StatusCode {
	case 400:
		var badRequestError BadRequestError
		err := ParseJSON(res, &badRequestError)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		return &RequestFailedError{res.StatusCode, badRequestError, req}
	case 401:
		return &RequestFailedError{res.StatusCode, errgo.New("unauthorized - you are not authorized to do this operation"), req}
	case 402:
		var paymentRequiredErr PaymentRequiredError
		err := ParseJSON(res, &paymentRequiredErr)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		return &RequestFailedError{res.StatusCode, paymentRequiredErr, req}
	case 404:
		var notFoundErr NotFoundError
		err := ParseJSON(res, &notFoundErr)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		return &RequestFailedError{res.StatusCode, notFoundErr, req}
	case 422:
		var unprocessableError UnprocessableEntity
		err := ParseJSON(res, &unprocessableError)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		return &RequestFailedError{res.StatusCode, unprocessableError, req}
	case 500:
		return &RequestFailedError{res.StatusCode, errgo.New("server internal error - our team has been notified"), req}
	case 503:
		return &RequestFailedError{res.StatusCode, fmt.Errorf("upstream provider returned an error, please retry later"), req}
	default:
		return &RequestFailedError{res.StatusCode, fmt.Errorf("invalid status from server: %v", res.Status), req}
	}
}

func (err *RequestFailedError) Error() string {
	return err.APIError.Error()
}

func (err *RequestFailedError) String() string {
	return err.APIError.Error()
}

func IsRequestFailedError(err error) bool {
	_, ok := err.(*RequestFailedError)
	return ok
}
