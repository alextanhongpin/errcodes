package errcodes

import "github.com/alextanhongpin/errcodes/internal"

type ErrorTrace = internal.ErrorTrace

func WithStack(err error) error {
	return internal.New(err)
}

func Wrap(err error, cause string) error {
	return internal.Wrap(err, cause)
}

func Sprint(err error, reversed bool) string {
	return internal.Sprint(err, reversed)
}

func StackTrace(err error) ([]uintptr, map[uintptr]string) {
	return internal.StackTrace(err)
}
