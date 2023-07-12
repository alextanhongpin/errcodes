package stacktrace_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errcodes/stacktrace"
)

var ErrUserNotFound = errors.New("user not found")

func ExampleStackTraceWithStack() {
	err := findUser()

	fmt.Println(stacktrace.Sprint(err, false))
	fmt.Println()
	fmt.Println("Reversed:")
	fmt.Println(stacktrace.Sprint(err, true))

	// Output:
	// Error: user not found
	//     Origin is:
	//         at stacktrace_test.findUser (in examples_stacktrace_with_stack_test.go:36)
	//     Ends here:
	//         at stacktrace_test.ExampleStackTraceWithStack (in examples_stacktrace_with_stack_test.go:13)
	//
	// Reversed:
	// Error: user not found
	//     Ends here:
	//         at stacktrace_test.ExampleStackTraceWithStack (in examples_stacktrace_with_stack_test.go:13)
	//     Origin is:
	//         at stacktrace_test.findUser (in examples_stacktrace_with_stack_test.go:36)
}

func findUser() error {
	err := stacktrace.WithStack(ErrUserNotFound)
	return err
}
