package errcodes_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errcodes"
)

func ExampleUnwrapStackTrace() {
	err := errcodes.WithStack(errors.New("bad request"))

	// Unwrap using errors.As.
	var errTrace *errcodes.ErrorTrace
	fmt.Println(errors.As(err, &errTrace))

	// Returns the raw unfiltered stacktrace.
	fmt.Println(len(errTrace.StackTrace()))

	// Returned the deduped stacktrace, together with cause annotation at
	// specific PCs.
	pcs, causes := errcodes.StackTrace(err)
	fmt.Println(len(pcs))
	fmt.Println(len(causes))
	fmt.Println(errcodes.Sprint(err, false))

	// Output:
	// true
	// 7
	// 1
	// 0
	// Error: bad request
	//     Origin is:
	//         at errcodes_test.ExampleUnwrapStackTrace (in examples_unwrap_stacktrace_test.go:11)
}
