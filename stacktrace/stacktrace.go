package stacktrace

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
)

var gopath = os.Getenv("GOPATH")

func Sprint(err error) string {
	var res []string
	res = append(res, fmt.Sprintf("Error: %s", err.Error()))

	var e *errorStack
	if errors.As(err, &e) {
		// Include the current line into the stack trace.
		e.pcs = mergePCs(e.pcs, caller(2))
		res = append(res, e.String())
	}

	return strings.Join(res, "\n")
}

// WithStack adds the stack trace at the current line.
func WithStack(err error) error {
	var errStack *errorStack
	if errors.As(err, &errStack) {
		return err
	}

	return newErrorStack(err, "")
}

func WithCause(err error, msg string) error {
	var e *errorStack
	if errors.As(err, &e) {
		c := newErrorStack(err, msg)

		// Combine both PCs.
		c.pcs = mergePCs(c.pcs, e.StackTrace()...)

		// Copy kv from old error to new error.
		for k, v := range e.cause {
			c.cause[k] = v
		}

		return c
	}

	return newErrorStack(err, msg)
}

type errorStack struct {
	err   error
	pcs   []uintptr
	cause map[uintptr]string
}

func newErrorStack(err error, msg string) *errorStack {
	// skip [runtime, callers, newErrorStack, With*]
	pcs := callers(4)
	cause := make(map[uintptr]string)
	if msg != "" {
		// Use the function PC when annotating cause.
		// All lines from the same function will be grouped later.
		pc := frame(pcs[0]).Entry
		cause[pc] = msg
	}

	return &errorStack{
		err:   err,
		cause: cause,
		pcs:   pcs,
	}
}

func (e *errorStack) StackTrace() []uintptr {
	pcs := make([]uintptr, len(e.pcs))
	copy(pcs, e.pcs)
	return pcs
}

func (e *errorStack) Unwrap() error {
	return e.err
}

func (e *errorStack) Error() string {
	return e.err.Error()
}

func (e *errorStack) String() string {
	pcs := e.StackTrace()

	// Top-down, from the main program to the root cause.
	reverse(pcs)

	type funcStack struct {
		frames []runtime.Frame
		cause  string
	}

	var orderedPCs []uintptr

	groupByFuncPC := make(map[uintptr]funcStack)

	frames := runtime.CallersFrames(pcs)

	for {
		f, more := frames.Next()

		if strings.HasPrefix(f.Function, "runtime") {
			continue
		}
		if f.Function == "" {
			break
		}

		// We use the function PC as a reference.
		// We want to group all the lines by the function PC,
		// then sort them by line number in descending order
		// (bottom-up) so that we can trace back errors.
		pc := f.Entry

		var cause string
		if msg, ok := e.cause[pc]; ok && len(msg) > 0 {
			cause = msg
		}

		// Some lines may be out of order.
		// However, due to how the PC are always incrementing,
		// if the PC is lower, it means it is called first.
		// We will group other PCs that comes later by the same function PC.
		if _, ok := groupByFuncPC[pc]; !ok {
			groupByFuncPC[f.Entry] = funcStack{}
			orderedPCs = append(orderedPCs, pc)
		}

		g := groupByFuncPC[pc]
		g.cause = cause
		g.frames = append(g.frames, f)
		groupByFuncPC[pc] = g

		if !more {
			break
		}
	}

	var out []string

	exists := make(map[string]bool)
	for i := 0; i < len(orderedPCs); i++ {
		pc := orderedPCs[i]
		g := groupByFuncPC[pc]

		fs := g.frames

		// Order the frames from same function in descending order.
		// This will appear bottom-up, so it is easier to trace the error.
		sort.Slice(fs, func(i, j int) bool {
			return fs[i].Line > fs[j].Line
		})

		// Annotate at function level.
		if g.cause != "" {
			out = append(out, fmt.Sprintf("    Caused by: %s", g.cause))
		}

		for j := range fs {
			f := fs[j]
			s := formatFrame(f)
			if exists[s] {
				continue
			}

			out = append(out, fmt.Sprintf("\t%s", s))
			exists[s] = true
		}

		out = append(out, "")
	}

	return strings.Join(out, "\n")
}

func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func formatFrame(frame runtime.Frame) string {
	return fmt.Sprintf("at %s (in %s:%d)",
		frame.Function,
		prettyFile(frame.File),
		frame.Line,
	)
}

func typeName(v any) (res string) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		res += "*"
		t = t.Elem()
	}

	if p := t.PkgPath(); p != "" {
		res += p
		res += "."
	}
	res += t.Name()
	return
}

func frame(pc uintptr) runtime.Frame {
	f, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	return f
}

func caller(skip int) uintptr {
	pc, _, _, _ := runtime.Caller(skip)
	return pc
}

func callers(skip int) []uintptr {
	const depth = 64
	var pc [depth]uintptr
	n := runtime.Callers(skip, pc[:])
	if n == 0 {
		return nil
	}

	var pcs = pc[:n]
	return pcs
}

func prettyFile(f string) string {
	if len(gopath) == 0 {
		return f
	}

	// TODO: also split by bitbucket, github.com, gopkg, golang.org?

	parts := strings.Split(f, fmt.Sprintf("%s/src", gopath))
	part := parts[len(parts)-1]

	return strings.TrimPrefix(part, "/")
}

func mergePCs(x []uintptr, y ...uintptr) []uintptr {
	z := append(x, y...)

	sort.Slice(z, func(i, j int) bool {
		return z[i] < z[j]
	})

	return slices.Compact(z)
}
