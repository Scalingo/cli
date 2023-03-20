package errors

import (
	"reflect"

	"gopkg.in/errgo.v1"
)

// IsRootCause return true if the cause of the given error is the same type as
// mytype.
// This function takes the cause of an error if the errors stack has been
// wrapped with errors.Wrapf or errgo.Notef or errgo.NoteMask or errgo.Mask.
//
// Example:
//
//	errors.IsRootCause(err, &ValidationErrors{})
func IsRootCause(err error, mytype interface{}) bool {
	t := reflect.TypeOf(mytype)
	errCause := errorCause(err)
	errRoot := errgoRoot(err)
	return reflect.TypeOf(errCause) == t || reflect.TypeOf(errRoot) == t
}

// RootCause returns the cause of an errors stack, whatever the method they used
// to be stacked: either errgo.Notef or errors.Wrapf.
func RootCause(err error) error {
	errCause := errorCause(err)
	if errCause == nil {
		errCause = errgoRoot(err)
	}
	return errCause
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
