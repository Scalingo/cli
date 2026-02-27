package errors

import "errors"

// Join combines multiple errors into a single error.
// It returns nil if all errors are nil.
// It uses errors.Join from the standard library}
func Join(errs ...error) error {
	return errors.Join(errs...)
}
