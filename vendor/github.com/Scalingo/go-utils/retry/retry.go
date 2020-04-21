package retry

import (
	"context"
	"fmt"
	"math"
	"time"
)

type RetryErrorScope string

const (
	MaxDurationScope RetryErrorScope = "max-duration"
	ContextScope     RetryErrorScope = "context"
)

type RetryError struct {
	Scope   RetryErrorScope
	Err     error
	LastErr error
}

func (err RetryError) Error() string {
	return fmt.Sprintf("retry error (%v): %v, last error %v", err.Scope, err.Err, err.LastErr)
}

// RetryCancelError is a error wrapping type that the user of a Retry should
// use to cancel retry operations before the end of maxAttempts/maxDuration
// conditions
type RetryCancelError struct {
	error
}

func NewRetryCancelError(err error) RetryCancelError {
	return RetryCancelError{error: err}
}

func (err RetryCancelError) Error() string {
	return err.error.Error()
}

type Retryable func(ctx context.Context) error

type ErrorCallback func(ctx context.Context, err error, currentAttempt, maxAttempts int)

type Retry interface {
	Do(ctx context.Context, method Retryable) error
}

type Retryer struct {
	waitDuration   time.Duration
	maxDuration    time.Duration
	maxAttempts    int
	errorCallbacks []ErrorCallback
}

type RetryerOptsFunc func(r *Retryer)

func WithWaitDuration(duration time.Duration) RetryerOptsFunc {
	return func(r *Retryer) {
		r.waitDuration = duration
	}
}

func WithMaxAttempts(maxAttempts int) RetryerOptsFunc {
	return func(r *Retryer) {
		r.maxAttempts = maxAttempts
	}
}

func WithMaxDuration(duration time.Duration) RetryerOptsFunc {
	return func(r *Retryer) {
		r.maxDuration = duration
	}
}

func WithoutMaxAttempts() RetryerOptsFunc {
	return func(r *Retryer) {
		r.maxAttempts = math.MaxInt32
	}
}

func WithErrorCallback(c ErrorCallback) RetryerOptsFunc {
	return func(r *Retryer) {
		r.errorCallbacks = append(r.errorCallbacks, c)
	}
}

func New(opts ...RetryerOptsFunc) Retryer {
	r := &Retryer{
		waitDuration:   10 * time.Second,
		maxAttempts:    5,
		errorCallbacks: make([]ErrorCallback, 0),
	}

	for _, opt := range opts {
		opt(r)
	}

	return *r
}

// Do execute method following rules of the Retry struct
// Two timeouts co-exist:
// * The one given as param of 'method': can be the scope of the current
// http.Request for instance
// * The one defined with the option WithMaxDuration, which would cancel the
// retry loop if it has expired.
func (r Retryer) Do(ctx context.Context, method Retryable) error {
	timeoutCtx := context.Background()
	if r.maxDuration != 0 {
		var cancel func()
		timeoutCtx, cancel = context.WithTimeout(timeoutCtx, r.maxDuration)
		defer cancel()
	}

	var err error
	for i := 0; i < r.maxAttempts; i++ {
		err = method(ctx)
		if err == nil {
			return nil
		}
		if rerr, ok := err.(RetryCancelError); ok {
			return rerr.error
		}

		for _, c := range r.errorCallbacks {
			c(ctx, err, i, r.maxAttempts)
		}

		timer := time.NewTimer(r.waitDuration)
		select {
		case <-timer.C:
		case <-timeoutCtx.Done():
			return RetryError{
				Scope:   MaxDurationScope,
				Err:     timeoutCtx.Err(),
				LastErr: err,
			}
		case <-ctx.Done():
			return RetryError{
				Scope:   ContextScope,
				Err:     ctx.Err(),
				LastErr: err,
			}
		}
	}
	return err
}
