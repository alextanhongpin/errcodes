package errcodes

import (
	"errors"
)

func Details(err error) []any {
	var e *errorDetails
	if errors.As(err, &e) {
		return e.details
	}

	return nil
}

func DetailsAs[T any](err error) ([]T, bool) {
	var e *errorDetails
	if errors.As(err, &e) {
		details := make([]T, len(e.details))
		for i := 0; i < len(e.details); i++ {
			t, ok := As[T](e.details[i])
			if !ok {
				return nil, false
			}
			details[i] = t
		}

		return details, true
	}

	return nil, false
}

func NewDetails(err error, details ...any) error {
	return &errorDetails{
		err:     err,
		details: details,
	}
}

type errorDetails struct {
	err     error
	details []any
}

func (e *errorDetails) Error() string {
	return e.err.Error()
}

func (e *errorDetails) Unwrap() error {
	return e.err
}

func As[T any](v any) (T, bool) {
	t, ok := v.(T)
	return t, ok
}
