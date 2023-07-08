package stacktrace

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

const head = "Error:"
const origin = "Origin is:"
const end = "Ends here:"

// StackTrace returns the stacktrace from nested errors after removing
// duplicated program counters.
func StackTrace(err error) []uintptr {
	var stack []uintptr

	seen := make(map[runtime.Frame]bool)

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

		var ordered []uintptr
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
			// Instead of f.PC + 1, easier to point to original stack.
			ordered = append(ordered, r.stack[i])
		}

		// Stack is ordered from bottom-up.
		// Reverse it.
		reverse(ordered)

		stack = append(stack, ordered...)
		err = r.Unwrap()
	}

	// Return in the order as what the original
	// runtime.Callers will return, which is
	// from the error origin to main.
	reverse(stack)

	return stack
}

func SprintCaller(err error, skip int, reversed ...bool) string {
	return sprintCaller(err, len(reversed) > 0, 1+skip)
}

// Sprint prints a readable stacktrace together with the cause.
func Sprint(err error, reversed ...bool) string {
	return sprintCaller(err, len(reversed) > 0, 1)
}

func Wrap(err error, cause string) error {
	return wrapCaller(err, cause, 1)
}

func WrapCaller(err error, cause string, skip int) error {
	return wrapCaller(err, cause, 1+skip)
}

func New(err error) error {
	return newCaller(err, 1)
}

func NewCaller(err error, skip int) error {
	return newCaller(err, 1+skip)
}

type errorStack struct {
	err   error
	stack []uintptr
}

func newErrorStack(err error, skip int) *errorStack {
	skip++

	return &errorStack{
		err: err,
		// skip [newErrorStack]
		stack: callers(skip),
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

// for each new call, add 1 to skip
func newErrorCause(err error, cause string, skip int) *errorCause {
	skip++

	if cause != "" {
		cause = fmt.Sprintf("Caused by: %s", cause)
	}
	return &errorCause{
		err: err,
		// skip [newErrorCause]
		stack: callers(skip),
		cause: cause,
	}
}

func (e *errorCause) Error() string {
	return e.err.Error()
}

func (e *errorCause) Unwrap() error {
	return e.err
}

func sprintCaller(err error, reversed bool, skip int) string {
	skip++

	var res []string

	header := fmt.Sprintf("%s %s", head, err)

	seen := make(map[runtime.Frame]bool)
	errC := newErrorCause(err, "", skip)
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

		// By default, it is ordered from inner (origin of error) to outer main
		// program.
		var ordered []runtime.Frame
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
			ordered = append(ordered, f)
		}

		var j int
		if reversed {
			reverse(ordered)
			j = len(ordered) - 1
		} else {
			j = 0
		}

		var tmp []string
		for i, f := range ordered {
			var s string
			if i == j && r.cause != "" {
				s = fmt.Sprintf("    %s\n", r.cause)
			}
			s = fmt.Sprintf("%s        %s", s, FormatFrame(f))
			tmp = append(tmp, s)
		}

		if !reversed {
			reverse(tmp)
		}

		res = append(res, tmp...)
		err = r.Unwrap()
	}

	if reversed {
		res = append([]string{header}, res...)
	} else {
		res = append(res, header)
		reverse(res)
	}

	return strings.Join(res, "\n")
}

func wrapCaller(err error, cause string, skip int) error {
	skip++
	if !isErrorStack(err) {
		errS := newErrorStack(err, skip)
		errC := newErrorCause(errS, cause, skip)
		errC.cause = fmt.Sprintf("%s %s", origin, cause)
		return errC
	}

	return newErrorCause(err, cause, skip)
}

func newCaller(err error, skip int) error {
	skip++
	if isErrorStack(err) {
		return newErrorCause(err, "", skip)
	}

	errC := newErrorCause(newErrorStack(err, skip), "", skip)
	errC.cause = origin
	return errC
}
