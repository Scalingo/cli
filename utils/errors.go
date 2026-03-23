package utils

import (
	"github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

func IsRegionDisabledError(err error) bool {
	var reqerr *http.RequestFailedError
	if !errors.As(err, &reqerr) || reqerr.Code != 403 {
		return false
	}
	httperr, ok := reqerr.APIError.(http.ForbiddenError)
	return ok && httperr.Code == "region_disabled"
}
