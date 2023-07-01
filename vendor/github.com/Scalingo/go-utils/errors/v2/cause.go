package errors

import (
	"errors"
	"reflect"

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

// RootCause returns the cause of an errors stack, whatever the method they used
// to be stacked: either errgo.Notef or errors.Wrapf.
//
// Deprecated: Use `Is(err, expectedErr)` instead of `if RootCause(err) == expectedErr` to match go standard libraries practices
func RootCause(err error) error {
	errCause := errorCause(err)
	if errCause == nil {
		errCause = errgoRoot(err)
	}
	return errCause
}

// IsRootCause return true if the cause of the given error is the same type as
// mytype.
// This function takes the cause of an error if the errors stack has been
// wrapped with errors.Wrapf or errgo.Notef or errgo.NoteMask or errgo.Mask.
//
// Example:
//
//	errors.IsRootCause(err, &ValidationErrors{})
//
// Deprecated: Use `As(err, mytype)` instead to match go standard libraries practices
func IsRootCause(err error, mytype interface{}) bool {
	t := reflect.TypeOf(mytype)
	errCause := errorCause(err)
	errRoot := errgoRoot(err)
	return reflect.TypeOf(errCause) == t || reflect.TypeOf(errRoot) == t
}

// UnwrapError tries to unwrap `err`. It unwraps any causer type, errgo and ErrCtx errors.
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

	type causer interface {
		Cause() error
	}

	// if err is type of `ErrCtx` unwrap it by getting errCtx.err
	if ctxerr, ok := err.(ErrCtx); ok {
		return ctxerr.err
	}

	// Check if the err is type of `*errgo.Err` to be able to call `Underlying()`
	// method. Both `*errgo.Err` and `*errors.Err` are implementing a causer interface.
	// Cause() method from errgo skip all underlying errors, so we may skip a context between.
	// So the order matter, we need to call `Cause()` after `Underlying()`.
	if errgoErr, ok := err.(*errgo.Err); ok {
		return errgoErr.Underlying()
	}

	if cause, ok := err.(causer); ok {
		return cause.Cause()
	}
	return nil
}
