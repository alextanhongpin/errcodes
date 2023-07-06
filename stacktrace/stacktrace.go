package main

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func bar() error {
	return WithStack(errors.New("bad"))
}

func foo() error {
	err := bar()

	err = WithCause(err, "hello")

	return err
}

func main() {
	err := foo()

	err = fmt.Errorf("%w: walalo", err)

	fmt.Println(Sprint(err))
}

func WithStack(err error) error {
	var errStack *errorStack
	if errors.As(err, &errStack) {
		return err
	}

	return newErrorStack(err)
}

func WithCause(err error, msg string) error {
	var errStack *errorStack
	if !errors.As(err, &errStack) {
		errStack = newErrorStack(err)
		fmt.Println("is not error stack", err, msg)
	}
	fmt.Println("is error stack", err, msg)

	c := errStack.clone()
	// Skips [caller, WithCause].
	pc := caller(2)
	frame, _ := runtime.CallersFrames([]uintptr{pc - 1}).Next()
	c.cause[formatFrame(frame)] = msg
	c.stack = append(c.stack, pc)

	return c
}

func Sprint(err error) string {
	var errStack *errorStack
	if !errors.As(err, &errStack) {
		return err.Error()
	}

	// Skips [caller, Sprint].
	pc := caller(2)
	frame, _ := runtime.CallersFrames([]uintptr{pc - 1}).Next()

	c := errStack.clone()
	c.err = err
	//c.stack = append([]uintptr{frame.PC}, c.stack...)
	c.stack = append(c.stack, frame.PC)
	return c.String()
	//return errStack.String()
}

type errorStack struct {
	err   error
	stack []uintptr
	cause map[string]string
}

func newErrorStack(err error) *errorStack {
	return &errorStack{
		err: err,
		// Skip [runtime, callers, and newErrorStack, and WithStack].
		stack: callers(4),
		cause: make(map[string]string),
	}
}

func (e *errorStack) clone() *errorStack {
	c := &errorStack{
		err:   e.err,
		stack: make([]uintptr, len(e.stack)),
		cause: make(map[string]string),
	}
	copy(c.stack, e.stack)
	for k, v := range e.cause {
		c.cause[k] = v
	}

	return c
}

func (e *errorStack) StackTrace() []uintptr {
	return e.stack
}

func (e *errorStack) Error() string {
	return e.err.Error()
}

func (err *errorStack) String() string {
	stack := err.stack
	cause := err.cause

	s := make([]uintptr, len(stack))
	copy(s, stack)
	reverse(s)

	var sb strings.Builder
	sb.WriteString(typeName(err.err) + ": " + err.err.Error())
	sb.WriteRune('\n')
	sb.WriteRune('\n')

	frames := runtime.CallersFrames(s)
	for {
		frame, more := frames.Next()

		if strings.HasPrefix(frame.Function, "runtime") {
			continue
		}

		if frame.Function == "" {
			continue
		}

		f := formatFrame(frame)
		if msg, ok := cause[f]; ok {
			sb.WriteString("    Caused by: ")
			sb.WriteString(msg)
			sb.WriteRune('\n')
		}
		sb.WriteRune('\t')
		sb.WriteString(f)
		sb.WriteRune('\n')

		if !more {
			break
		}
	}

	return sb.String()
}

func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func formatFrame(frame runtime.Frame) string {
	return fmt.Sprintf("at %s in %s:%d",
		frame.Function,
		frame.File,
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

func caller(skip int) uintptr {
	pc, _, _, _ := runtime.Caller(skip)
	return pc
}

func callers(skip int) []uintptr {
	const depth = 64
	var pc [depth]uintptr
	n := runtime.Callers(skip, pc[:])

	var pcs = pc[:n]
	return pcs
}
