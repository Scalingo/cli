package scalingo

import (
	"strconv"

	"github.com/Scalingo/go-utils/pagination"
)

func paginationRequestToMap(paginationReq pagination.Request) map[string]string {
	return map[string]string{
		"page":     strconv.Itoa(paginationReq.Page),
		"per_page": strconv.Itoa(paginationReq.PerPage),
	}
}
