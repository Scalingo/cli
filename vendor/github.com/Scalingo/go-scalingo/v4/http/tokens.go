package http

type TokenGenerator interface {
	GetAccessToken() (string, error)
}
