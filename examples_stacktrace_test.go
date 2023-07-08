package errcodes_test

import (
	"fmt"

	"github.com/alextanhongpin/errcodes"
)

var ErrUserNotFound = errcodes.New(errcodes.NotFound, "user_not_found", "The user does not exists or may have been deleted")

func ExampleStacktrace() {
	err := findUser()
	fmt.Println(errcodes.Sprint(err))
	// Output:
	// Error: The user does not exists or may have been deleted
	//     Ends here:
	//         at errcodes_test.ExampleStacktrace (in examples_stacktrace_test.go:13)
	//         at errcodes_test.ExampleStacktrace (in examples_stacktrace_test.go:12)
	//     Origin is:
	//         at errcodes_test.findUser (in examples_stacktrace_test.go:24)
}

func findUser() error {
	err := errcodes.WithStack(ErrUserNotFound)
	return err
}
