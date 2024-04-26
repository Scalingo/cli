package utils

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v7/http"
	errors "github.com/Scalingo/go-utils/errors/v2"
)

func IsRegionDisabledError(err error) bool {
	reqerr, ok := errors.RootCause(err).(*http.RequestFailedError)
	if !ok || reqerr.Code != 403 {
		return false
	}
	httperr, ok := reqerr.APIError.(http.ForbiddenError)
	return ok && httperr.Code == "region_disabled"
}

func WrapError(err error, wrappingMessage string) error {
	return errgo.Notef(err, wrappingMessage)
}
