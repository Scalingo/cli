package http

import "context"

type TokenGenerator interface {
	GetAccessToken(context.Context) (string, error)
}
