package api

type InternalError struct {
	Error string `json:"error"`
}

type BadRequestError struct {
	Errors map[string][]string `json:"errors"`
}
