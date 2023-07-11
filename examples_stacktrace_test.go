package errcodes_test

import (
	"fmt"

	"github.com/alextanhongpin/errcodes"
)

var ErrUserNotFound = errcodes.New(errcodes.NotFound, "user_not_found", "The user does not exists or may have been deleted")

func ExampleStackTrace() {
	err := findUser()
	fmt.Println(errcodes.Sprint(err, false))
	fmt.Println()
	fmt.Println("Reversed:")
	fmt.Println(errcodes.Sprint(err, true))

	// Output:
	// Error: The user does not exists or may have been deleted
	//     Origin is:
	//         at errcodes_test.findUser (in examples_stacktrace_test.go:34)
	//     Ends here:
	//         at errcodes_test.ExampleStackTrace (in examples_stacktrace_test.go:12)
	//
	// Reversed:
	// Error: The user does not exists or may have been deleted
	//     Ends here:
	//         at errcodes_test.ExampleStackTrace (in examples_stacktrace_test.go:12)
	//     Origin is:
	//         at errcodes_test.findUser (in examples_stacktrace_test.go:34)
}

func findUser() error {
	err := errcodes.WithStack(ErrUserNotFound)
	return err
}
