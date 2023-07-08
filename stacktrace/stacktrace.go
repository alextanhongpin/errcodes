package stacktrace

import (
	"errors"
	"runtime"
	"strings"
)

// Flatten flattens the stacktrace from nested errors and remove duplicates
// program counters.
// Use this before sending the stacktrace information to monitoring tools like
// Sentry etc.
// Flatten removes the Cause from the *errorDetails.
func Flatten(err error) error {
	var stack []uintptr

	seen := make(map[runtime.Frame]bool)
	err = newErrorStack(err, "")

	for err != nil {
		var r *errorStack
		if !errors.As(err, &r) {
			break
		}

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

			stack = append(stack, f.PC)
		}

		err = r.Unwrap()
	}

	return &errorStack{
		err:   err,
		stack: stack,
	}
}

// Sprint prints a readable stacktrace together with the cause.
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

		var rev []runtime.Frame
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

			rev = append(rev, f)
		}

		reverse(rev)

		for i, f := range rev {
			if i == len(rev)-1 && r.cause != "" {
				sb.WriteString("    Caused by: ")
				sb.WriteString(r.cause)
				sb.WriteRune('\n')
			}
			sb.WriteRune('\t')
			sb.WriteString(formatFrame(f))
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

// StackTrace returns the stacktrace of the current error.
func (r *errorStack) StackTrace() []uintptr {
	return r.stack
}

func (r *errorStack) Error() string {
	return r.err.Error()
}

func (r *errorStack) Unwrap() error {
	return r.err
}
