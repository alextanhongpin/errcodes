package stacktrace_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errcodes/stacktrace"
)

func foo() error {
	err := errors.New("foo")
	return stacktrace.New(err)
}

func bar() error {
	err := foo()
	return stacktrace.Wrap(err, "bar")
}

func ExampleStackTrace() {
	fmt.Println(stacktrace.Sprint(bar()))
	// Output:
	// Error: foo
	// 	at stacktrace_test.ExampleStackTrace (in stacktrace_test.go:21)
	//     Caused by: bar
	// 	at stacktrace_test.bar (in stacktrace_test.go:17)
	// 	at stacktrace_test.bar (in stacktrace_test.go:16)
	//     Caused by: foo
	// 	at stacktrace_test.foo (in stacktrace_test.go:12)
}
