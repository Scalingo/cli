package errors

import "gopkg.in/errgo.v1"

func ErrgoRoot(err error) error {
	for {
		e, ok := err.(*errgo.Err)
		if !ok {
			return err
		}
		if e.Underlying() == nil {
			return err
		}
		err = e.Underlying()
	}
}
