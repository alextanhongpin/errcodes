package errcodes_test

import (
	"fmt"

	"github.com/alextanhongpin/errcodes"
)

var ErrUserNotFound = errcodes.New(errcodes.NotFound, "user_not_found", "The user does not exists or may have been deleted")

func ExampleStacktrace() {
	err := findUser()
	fmt.Println(errcodes.Sprint(err))
	fmt.Println()
	fmt.Println("Reversed:")
	fmt.Println(errcodes.Sprint(err, true))

	// Output:
	// Error: The user does not exists or may have been deleted
	//     Origin is:
	//         at errcodes_test.findUser (in examples_stacktrace_test.go:36)
	//         at errcodes_test.ExampleStacktrace (in examples_stacktrace_test.go:12)
	//     Ends here:
	//         at errcodes_test.ExampleStacktrace (in examples_stacktrace_test.go:13)
	//
	// Reversed:
	// Error: The user does not exists or may have been deleted
	//     Ends here:
	//         at errcodes_test.ExampleStacktrace (in examples_stacktrace_test.go:16)
	//         at errcodes_test.ExampleStacktrace (in examples_stacktrace_test.go:12)
	//     Origin is:
	//         at errcodes_test.findUser (in examples_stacktrace_test.go:36)
}

func findUser() error {
	err := errcodes.WithStack(ErrUserNotFound)
	return err
}
