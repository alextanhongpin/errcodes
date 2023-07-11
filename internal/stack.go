package internal

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

const depth = 32

const indent = "    "

type node int

const (
	none node = iota // Don't expose stacktrace.
	root             // There can only be one root with stacktrace.
	leaf             // There can be multiple leaf with stacktrace.
)

const head = "Origin is:"
const tail = "Ends here:"
const body = "Caused by:"

func New(err error) error {
	return &ErrorTrace{
		node:  root,
		err:   err,
		stack: callers(2), // Skips [New, caller]
	}
}

func Wrap(err error, cause string) error {
	if err == nil {
		return nil
	}

	// Skips [Wrap, caller]
	return wrap(err, cause, 2)
}

func Sprint(err error, reversed bool) string {
	// Skips [Sprint, Caller]
	return sprint(err, reversed, 2)
}

func StackTrace(err error) ([]uintptr, map[uintptr]string) {
	return unwrap(err)
}

func wrap(err error, cause string, skip int) *ErrorTrace {
	if err == nil {
		return nil
	}

	// Skip [wrap]
	skip = skip + 1

	var t *ErrorTrace
	if errors.As(err, &t) {
		// Can happen when for loop wrapping the same error
		// at the same line of code.
		if err == t && t.cause == cause {
			return t
		}

		errT := extractLastRootOrLeafError(err)
		// The trace is deep, add a new leaf node.
		if len(errT.StackTrace()) == depth {
			return &ErrorTrace{
				node:  leaf,
				err:   err,
				stack: callers(skip),
				cause: cause,
			}
		}

		// It's not that deep, mark it first.
		return &ErrorTrace{
			node:  none,
			err:   err,
			stack: callers(skip),
			cause: cause,
		}
	}

	return &ErrorTrace{
		node:  root,
		err:   err,
		stack: callers(skip),
		cause: cause,
	}
}

type ErrorTrace struct {
	node  node
	err   error
	stack []uintptr
	cause string
}

func (e *ErrorTrace) StackTrace() []uintptr {
	// Only expose at the root and the leaf node.
	if e.node != none {
		return e.stack
	}

	return nil
}

func (e *ErrorTrace) Error() string {
	return e.err.Error()
}

func (e *ErrorTrace) Unwrap() error {
	return e.err
}

func extractLastRootOrLeafError(err error) *ErrorTrace {
	for {
		var t *ErrorTrace
		if !errors.As(err, &t) {
			break
		}

		if t.node == root || t.node == leaf {
			return t
		}

		err = t.Unwrap()
	}

	panic("no root or leaf found")
}

func sprint(err error, reversed bool, skip int) string {
	if err == nil {
		return ""
	}
	// Skip [sprint]
	skip = skip + 1

	var sb strings.Builder

	sb.WriteString("Error:")
	sb.WriteRune(' ')
	sb.WriteString(err.Error())
	sb.WriteRune('\n')

	pcs, cause := unwrap(err)
	pcs, cause = prettyCause(pcs, cause)
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
		sb.WriteString(formatFrame(frame))
		if !more {
			break
		}

		sb.WriteRune('\n')
	}

	return sb.String()
}

func unwrap(err error) ([]uintptr, map[uintptr]string) {
	if err == nil {
		return nil, nil
	}

	var pcs []uintptr
	cause := make(map[uintptr]string)
	seen := make(map[runtime.Frame]bool)

	for err != nil {
		var t *ErrorTrace
		if !errors.As(err, &t) {
			break
		}

		var ordered []uintptr
		for _, f := range cleanFrames(t.stack) {
			key := runtime.Frame{
				File:     f.File,
				Function: f.Function,
				Line:     f.Line,
			}
			if seen[key] {
				break
			}
			seen[key] = true
			// The runtime.CallersFrames PC =
			// runtime.callers(skip) PC - 1
			ordered = append(ordered, f.PC+1)
		}

		// The first frame indicates the cause.
		if len(ordered) > 0 && len(t.cause) > 0 {
			cause[ordered[0]] = t.cause
		}

		// Stack is ordered from bottom-up.
		// Reverse it so that it goes top-down.
		reverse(ordered)

		pcs = append(pcs, ordered...)
		err = t.Unwrap()
	}

	// Return in the order as what the original
	// runtime.callers will return, which is bottom-up.
	reverse(pcs)

	return pcs, cause
}

func cleanFrames(pcs []uintptr) []runtime.Frame {
	var res []runtime.Frame
	frames := runtime.CallersFrames(pcs)
	for {
		f, more := frames.Next()
		if !skipFrame(f) {
			res = append(res, f)
		}

		if !more {
			break
		}
	}

	return res
}

func skipFrame(f runtime.Frame) bool {
	// Skip empty function.
	return f.Function == "" ||
		// Skip runtime and testing package.
		strings.HasPrefix(f.Function, "runtime") ||
		strings.HasPrefix(f.Function, "testing") ||

		// Skip files with underscore.
		// e.g. _testmain.go
		strings.HasPrefix(f.File, "_")
}

func formatFrame(frame runtime.Frame) string {
	return fmt.Sprintf("at %s (in %s:%d)",
		prettyFunction(frame.Function),
		prettyFile(frame.File),
		frame.Line,
	)
}

func prettyFile(f string) string {
	wd, err := os.Getwd()
	if err != nil {
		return f
	}

	f = strings.TrimPrefix(f, wd)
	return strings.TrimPrefix(f, "/")
}

func prettyFunction(f string) string {
	_, file := path.Split(f)
	return file
}

func prettyCause(pcs []uintptr, cause map[uintptr]string) ([]uintptr, map[uintptr]string) {
	switch len(pcs) {
	case 0:
	case 1:
		pc := pcs[0]
		if msg, ok := cause[pc]; ok {
			cause[pc] = fmt.Sprintf("%s %s", head, msg)
		} else {
			cause[pc] = head
		}
	default:
		pc := pcs[0]
		if msg, ok := cause[pc]; ok {
			cause[pc] = fmt.Sprintf("%s %s", head, msg)
		} else {
			cause[pc] = head
		}

		for pc := range cause {
			if pc == pcs[0] || pc == pcs[len(pcs)-1] {
				continue
			}

			if msg, ok := cause[pc]; ok {
				cause[pc] = fmt.Sprintf("%s %s", body, msg)
			}
		}

		pc = pcs[len(pcs)-1]
		if msg, ok := cause[pc]; ok {
			cause[pc] = fmt.Sprintf("%s %s", tail, msg)
		} else {
			cause[pc] = tail
		}
	}
	return pcs, cause
}

func callers(skip int) []uintptr {
	var pc [depth]uintptr
	// skip [runtime.callers, callers]
	n := runtime.Callers(skip+2, pc[:])
	if n == 0 {
		return nil
	}

	var pcs = pc[:n]
	return pcs
}

func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
