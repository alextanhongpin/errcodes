package stacktrace

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

const indent = "    "
const head = "Error:"
const origin = "Origin is:"
const end = "Ends here:"

// Unwrap returns the stacktrace from nested errors after removing
// duplicated program counters.
func Unwrap(err error) ([]uintptr, map[uintptr]string) {
	return unwrap(err)
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

	errC := newErrorCause(err, "", skip)
	errC.cause = end
	err = errC

	var sb strings.Builder

	sb.WriteString(head)
	sb.WriteRune(' ')
	sb.WriteString(err.Error())
	sb.WriteRune('\n')

	pcs, cause := unwrap(err)
	if reversed {
		reverse(pcs)
	}

	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()

		msg, ok := cause[frame.PC+1]
		if ok && msg != "" {
			sb.WriteString(indent)
			sb.WriteString(msg)
			sb.WriteRune('\n')
		}
		sb.WriteString(indent)
		sb.WriteString(indent)
		sb.WriteString(FormatFrame(frame))
		if !more {
			break
		}

		sb.WriteRune('\n')
	}

	return sb.String()
}

func wrapCaller(err error, cause string, skip int) error {
	skip++
	if isErrorStack(err) {
		return newErrorCause(err, cause, skip)
	}

	errS := newErrorStack(err, skip)
	errC := newErrorCause(errS, cause, skip)
	errC.cause = fmt.Sprintf("%s %s", origin, cause)
	return errC
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

func unwrap(err error) ([]uintptr, map[uintptr]string) {
	var stack []uintptr
	cause := make(map[uintptr]string)
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
		if len(ordered) > 0 {
			cause[ordered[0]] = r.cause
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

	return stack, cause
}
