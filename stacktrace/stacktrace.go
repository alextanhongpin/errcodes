package stacktrace

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

func Sprint(err error) string {
	var sb strings.Builder

	sb.WriteString("Error: ")
	sb.WriteString(err.Error())
	sb.WriteRune('\n')

	seen := make(map[runtime.Frame]bool)
	err = newErrorStack(err, "")

	for err != nil {
		var r *errorStack
		if !errors.As(err, &r) {
			break
		}

		var rev []string
		for _, f := range frames(r.stack) {
			fi := runtime.Frame{
				File:     f.File,
				Function: f.Function,
				Line:     f.Line,
			}

			if seen[fi] {
				break
			}
			seen[fi] = true

			rev = append(rev, fmt.Sprintf("\t%s", formatFrame(f)))
		}

		reverse(rev)

		for i, s := range rev {
			if i == len(rev)-1 && r.cause != "" {
				sb.WriteString("    Caused by: ")
				sb.WriteString(r.cause)
				sb.WriteRune('\n')
			}
			sb.WriteString(s)
			sb.WriteRune('\n')
		}

		err = r.Unwrap()
	}

	return sb.String()
}

func Wrap(err error, cause string) error {
	return newErrorStack(err, cause)
}

func New(err error) error {
	return newErrorStack(err, err.Error())
}

type errorStack struct {
	err   error
	stack []uintptr
	cause string
}

func newErrorStack(err error, cause string) error {
	return &errorStack{
		err: err,
		// skip [runtime, caller, newErrorStack, parent]
		stack: callers(4),
		cause: cause,
	}
}

func (r *errorStack) Error() string {
	return r.err.Error()
}

func (r *errorStack) Unwrap() error {
	return r.err
}
