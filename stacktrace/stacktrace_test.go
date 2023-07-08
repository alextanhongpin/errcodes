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

func baz() error {
	err := errors.New("baz")
	// 1 means skip baz from being printed.
	return stacktrace.WrapCaller(err, "bzz!", 1)
}

func ExampleSprint() {
	fmt.Println(stacktrace.Sprint(bar()))
	// Output:
	// Error: foo
	//     Origin is:
	//         at stacktrace_test.foo (in stacktrace_test.go:13)
	//         at stacktrace_test.bar (in stacktrace_test.go:17)
	//     Caused by: bar
	//         at stacktrace_test.bar (in stacktrace_test.go:18)
	//     Ends here:
	//         at stacktrace_test.ExampleSprint (in stacktrace_test.go:28)
}

func ExampleSprintReverse() {
	fmt.Println(stacktrace.Sprint(bar(), true))
	// Output:
	// Error: foo
	//     Ends here:
	//         at stacktrace_test.ExampleSprintReverse (in stacktrace_test.go:41)
	//     Caused by: bar
	//         at stacktrace_test.bar (in stacktrace_test.go:18)
	//         at stacktrace_test.bar (in stacktrace_test.go:17)
	//     Origin is:
	//         at stacktrace_test.foo (in stacktrace_test.go:13)
}

func ExampleFormat() {
	err := bar()
	stack := stacktrace.StackTrace(err)
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		fmt.Println(stacktrace.FormatFrame(frame))
		if !more {
			break
		}
	}
	// Output:
	// at stacktrace_test.foo (in stacktrace_test.go:13)
	// at stacktrace_test.bar (in stacktrace_test.go:17)
	// at stacktrace_test.bar (in stacktrace_test.go:18)
	// at stacktrace_test.ExampleFormat (in stacktrace_test.go:54)
}

func ExampleCaller() {
	fmt.Println(stacktrace.Sprint(baz()))
	// Output:
	// Error: baz
	//     Ends here:
	//         at stacktrace_test.ExampleCaller (in stacktrace_test.go:72)
}
