package utils

import (
	"github.com/Scalingo/go-scalingo/http"
	"github.com/Scalingo/go-utils/errors"
)

func IsRegionDisabledError(err error) bool {
	reqerr, ok := errors.ErrgoRoot(err).(*http.RequestFailedError)
	if !ok || reqerr.Code != 403 {
		return false
	}
	httperr, ok := reqerr.APIError.(http.ForbiddenError)
	return ok && httperr.Code == "region_disabled"
}
