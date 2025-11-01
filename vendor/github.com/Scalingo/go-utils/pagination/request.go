package pagination

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
func (p Request) QueryLimit() int32 {
	// Returns an int32 so that the value can be used as a query argument without typecasting
	return int32(p.PerPage)
}

// QueryOffset The offset value to use when querying the DB
func (p Request) QueryOffset() int32 {
	// Returns an int32 so that the value can be used as a query argument without typecasting
	return int32((p.Page - 1) * p.PerPage)
}
