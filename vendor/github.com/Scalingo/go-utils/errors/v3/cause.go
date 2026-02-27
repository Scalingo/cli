package errors

import (
	"errors"

	"gopkg.in/errgo.v1"
)

// Is checks if any error of the stack matches the error value expectedError
// API machting the standard library but allowing to wrap errors with ErrCtx + errgo or pkg/errors
func Is(receivedErr, expectedError error) bool {
	if errors.Is(receivedErr, expectedError) {
		return true
	}
	for receivedErr != nil {
		receivedErr = UnwrapError(receivedErr)
		if errors.Is(receivedErr, expectedError) {
			return true
		}
	}
	return false
}

// As checks if any error of the stack matches the expectedType
// API machting the standard library but allowing to wrap errors with ErrCtx + errgo or pkg/errors
func As(receivedErr error, expectedType any) bool {
	if errors.As(receivedErr, expectedType) {
		return true
	}
	for receivedErr != nil {
		receivedErr = UnwrapError(receivedErr)
		if errors.As(receivedErr, expectedType) {
			return true
		}
	}
	return false
}

// UnwrapError tries to unwrap `err`. It unwraps any causer type, errgo and errors
// implementing Unwrap() method.
// It returns nil if no err found. This provide the possibility to loop on UnwrapError
// by checking the return value.
// E.g.:
//
//	for unwrappedErr := err; unwrappedErr != nil; unwrappedErr = UnwrapError(unwrappedErr) {
//		...
//	}
func UnwrapError(err error) error {
	if err == nil {
		return nil
	}

	// This also match errCtx from this package.
	u, ok := err.(interface {
		Unwrap() error
	})
	if ok {
		return u.Unwrap()
	}

	// Check if the err is type of `*errgo.Err` to be able to call `Underlying()`
	// method. Both `*errgo.Err` and `*errors.Err` are implementing a causer interface.
	// Cause() method from errgo skip all underlying errors, so we may skip a context between.
	// So the order matter, we need to call `Cause()` after `Underlying()`.
	if errgoErr, ok := err.(*errgo.Err); ok {
		return errgoErr.Underlying()
	}

	c, ok := err.(interface {
		Cause() error
	})
	if ok {
		return c.Cause()
	}

	return nil
}
