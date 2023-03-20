package errors

func errorCause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		e, ok := err.(ErrCtx)
		if ok {
			err = e.err
		}

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
