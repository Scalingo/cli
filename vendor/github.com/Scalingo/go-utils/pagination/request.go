package pagination

import (
	"net/url"
	"strconv"
)

const defaultPerPage = 20

type Request struct {
	Page    int `log:"page"`     // page requested (default 1)
	PerPage int `log:"per_page"` // Number of items per page
}

func NewRequest(page, perPage int) Request {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = defaultPerPage
	}
	return Request{
		Page:    page,
		PerPage: perPage,
	}
}

// QueryLimit The limit value to use when querying the DB
func (r Request) QueryLimit() int32 {
	// Returns an int32 so that the value can be used as a query argument without typecasting
	return int32(r.PerPage)
}

// QueryOffset The offset value to use when querying the DB
func (r Request) QueryOffset() int32 {
	// Returns an int32 so that the value can be used as a query argument without typecasting
	return int32((r.Page - 1) * r.PerPage)
}

func (r Request) ToURLValues() url.Values {
	values := url.Values{}
	values.Add("page", strconv.Itoa(r.Page))
	values.Add("per_page", strconv.Itoa(r.PerPage))
	return values
}
