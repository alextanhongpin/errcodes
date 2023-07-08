package stacktrace_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errcodes/stacktrace"
)

func ExampleStackTrace() {
	err := errors.New("an unexpected error")
	err = stacktrace.New(err)

	fmt.Println(stacktrace.Sprint(err))
	// Output:
	// Error: an unexpected error
	// 	at main.main (in _testmain.go:49)
	// 	at stacktrace_test.ExampleStackTrace (in stacktrace_test.go:14)
	//     Caused by: an unexpected error
	// 	at stacktrace_test.ExampleStackTrace (in stacktrace_test.go:12)
}
