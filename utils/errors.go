package utils

import (
	"github.com/Scalingo/go-scalingo/v9/http"
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
