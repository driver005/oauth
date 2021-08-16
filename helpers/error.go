package helper

import "github.com/pkg/errors"

type StackTracer interface {
	StackTrace() errors.StackTrace
}

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok || cause.Cause() == nil {
			break
		}
		err = cause.Cause()
	}
	return err
}

func WithStack(err error) error {
	if e, ok := err.(StackTracer); ok && len(e.StackTrace()) > 0 {
		return err
	}

	return errors.WithStack(err)
}
