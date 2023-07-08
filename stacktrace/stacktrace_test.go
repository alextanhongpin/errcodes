package stacktrace_test

import (
	"errors"
	"fmt"
	"runtime"

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

func ExampleSprint() {
	fmt.Println(stacktrace.Sprint(bar()))
	// Output:
	// Error: foo
	//     Ends here:
	//         at stacktrace_test.ExampleSprint (in stacktrace_test.go:22)
	//     Caused by: bar
	//         at stacktrace_test.bar (in stacktrace_test.go:18)
	//         at stacktrace_test.bar (in stacktrace_test.go:17)
	//     Origin is:
	//         at stacktrace_test.foo (in stacktrace_test.go:13)
}

func ExampleFormat() {
	stack := stacktrace.StackTrace(bar())
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		fmt.Println(stacktrace.FormatFrame(frame))
		if !more {
			break
		}
	}
	// Output:
	// at stacktrace_test.ExampleFormat (in stacktrace_test.go:35)
	// at stacktrace_test.bar (in stacktrace_test.go:18)
	// at stacktrace_test.bar (in stacktrace_test.go:17)
	// at stacktrace_test.foo (in stacktrace_test.go:13)
}
