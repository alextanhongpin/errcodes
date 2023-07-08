package errcodes_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errcodes"
)

var ErrDuplicateEmail = errcodes.New(errcodes.Conflict, "email_duplicate", "The email address is not available")

func ExampleErrorAs() {
	var err *errcodes.Error
	if errors.As(ErrDuplicateEmail, &err) {
		fmt.Println("1.", true)
	}

	fmt.Println("2.", err.Kind())
	fmt.Println("3.", err.Code())
	fmt.Println("4.", err.Message())
	fmt.Println("5.", err.String())

	// Output:
	// 1. true
	// 2. conflict
	// 3. email_duplicate
	// 4. The email address is not available
	// 5. conflict/email_duplicate: The email address is not available
}
