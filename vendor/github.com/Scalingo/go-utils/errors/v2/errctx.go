package errors

import (
	"context"

	"github.com/pkg/errors"
	"gopkg.in/errgo.v1"
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

// Unwrap implements error management from the standard library
func (err ErrCtx) Unwrap() error {
	return err.err
}

func New(ctx context.Context, message string) error {
	return ErrCtx{ctx: ctx, err: errors.New(message)}
}

func Newf(ctx context.Context, format string, args ...interface{}) error {
	return Errorf(ctx, format, args...)
}

// Deprecated: Use `Wrap` or `Wrapf` instead of `Notef`. The library is able to unwrap mixed errors (wrapped with `errgo` or `github.com/pkg/errors`).
func Notef(ctx context.Context, err error, format string, args ...interface{}) error {
	return ErrCtx{ctx: ctx, err: errgo.Notef(err, format, args...)}
}

func Wrap(ctx context.Context, err error, message string) error {
	return ErrCtx{ctx: ctx, err: errors.Wrap(err, message)}
}

func Wrapf(ctx context.Context, err error, format string, args ...interface{}) error {
	return ErrCtx{ctx: ctx, err: errors.Wrapf(err, format, args...)}
}

func Errorf(ctx context.Context, format string, args ...interface{}) error {
	return ErrCtx{ctx: ctx, err: errors.Errorf(format, args...)}
}

// RootCtxOrFallback unwrap all wrapped errors from err to get the deepest context
// from ErrCtx errors. If there is no wrapped ErrCtx RootCtxOrFallback returns ctx from parameter.
func RootCtxOrFallback(ctx context.Context, err error) context.Context {
	var lastCtx context.Context

	type causer interface {
		Cause() error
	}

	// Unwrap each error to get the deepest context
	for err != nil {
		// First check if the err is type of `*errgo.Err` to be able to call `Underlying()`
		// method. Both `*errgo.Err` and `*errors.Err` are implementing a causer interface.
		// Cause() method from errgo skip all underlying errors, so we may skip a context between.
		// So the order matter, we need to call `Cause()` after `Underlying()`.
		errgoErr, ok := err.(*errgo.Err)
		if ok {
			err = errgoErr.Underlying()
			continue
		}

		cause, ok := err.(causer)
		if ok {
			err = cause.Cause()
			continue
		}

		// if err is type of `ErrCtx` unwrap it by getting errCtx.err
		ctxerr, ok := err.(ErrCtx)
		if ok {
			err = ctxerr.err
			lastCtx = ctxerr.Ctx()

			continue
		}

		break
	}

	if lastCtx == nil {
		return ctx
	}

	return lastCtx
}
