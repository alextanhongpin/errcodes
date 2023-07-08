package stacktrace

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

const indent = "    "
const head = "Error: "
const origin = "Origin is:"
const end = "Ends here:"

// StackTrace returns the stacktrace from nested errors after removing
// duplicated program counters.
func StackTrace(err error) []uintptr {
	var stack []uintptr

	seen := make(map[runtime.Frame]bool)

	err = newErrorCause(err, "")

	for err != nil {
		var r *errorCause
		if !errors.As(err, &r) {
			var s *errorStack
			if !errors.As(err, &s) {
				break
			}

			r = &errorCause{
				stack: s.stack,
			}
		}

		var rev []uintptr
		for i, f := range frames(r.stack) {
			fi := runtime.Frame{
				File:     f.File,
				Function: f.Function,
				Line:     f.Line,
			}

			if seen[fi] {
				break
			}
			seen[fi] = true
			// The PC obtained by the runtime.Callers vs those from
			// runtime.CallersFrames, frame.PC differ by 1.
			rev = append(rev, r.stack[i])
		}

		reverse(rev)

		stack = append(stack, rev...)

		err = r.Unwrap()
	}

	return stack
}

// Sprint prints a readable stacktrace together with the cause.
func Sprint(err error) string {
	var sb strings.Builder

	sb.WriteString(head)
	sb.WriteString(err.Error())
	sb.WriteRune('\n')

	seen := make(map[runtime.Frame]bool)
	errC := newErrorCause(err, "")
	errC.cause = end
	err = errC

	for err != nil {
		var r *errorCause
		if !errors.As(err, &r) {
			var s *errorStack
			if !errors.As(err, &s) {
				break
			}

			r = &errorCause{
				stack: s.stack,
			}
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
				sb.WriteString(indent)
				sb.WriteString(r.cause)
				sb.WriteRune('\n')
			}
			sb.WriteString(indent)
			sb.WriteString(indent)
			sb.WriteString(FormatFrame(f))
			sb.WriteRune('\n')
		}
		err = r.Unwrap()
	}

	return sb.String()
}

func Wrap(err error, cause string) error {
	if !isErrorStack(err) {
		err = newErrorStack(err)
	}

	return newErrorCause(err, cause)
}

func New(err error) error {
	if isErrorStack(err) {
		return newErrorCause(err, "")
	}

	errC := newErrorCause(newErrorStack(err), "")
	errC.cause = origin
	return errC
}

type errorStack struct {
	err   error
	stack []uintptr
}

func newErrorStack(err error) *errorStack {
	return &errorStack{
		err: err,
		// skip [runtime, caller, newErrorStack, parent]
		stack: callers(4),
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

func isErrorStack(err error) bool {
	var e *errorStack
	return errors.As(err, &e)
}

type errorCause struct {
	err   error
	stack []uintptr
	cause string
}

func newErrorCause(err error, cause string) *errorCause {
	if cause != "" {
		cause = fmt.Sprintf("Caused by: %s", cause)
	}
	return &errorCause{
		err: err,
		// skip [runtime, caller, newErrorStack, parent]
		stack: callers(4),
		cause: cause,
	}
}

func (e *errorCause) Error() string {
	return e.err.Error()
}

func (e *errorCause) Unwrap() error {
	return e.err
}
