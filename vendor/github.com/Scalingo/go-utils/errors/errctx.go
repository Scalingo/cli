package errors

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/pkg/errors"
)

type ErrCtx struct {
	ctx context.Context
	err error
}

func (err ErrCtx) Error() string {
	return err.err.Error()
}

func (err ErrCtx) Ctx() context.Context {
	return err.ctx
}

func Notef(ctx context.Context, err error, format string, args ...interface{}) error {
	return ErrCtx{ctx: ctx, err: errgo.Notef(err, format, args...)}
}

func Wrapf(ctx context.Context, err error, format string, args ...interface{}) error {
	return ErrCtx{ctx: ctx, err: errors.Wrapf(err, format, args...)}
}
