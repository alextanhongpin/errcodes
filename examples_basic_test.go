package errcodes_test

import (
	"fmt"

	"github.com/alextanhongpin/errcodes"
)

var ErrAccountExists = errcodes.New(errcodes.Exists, "account_exists", "The user account already exists")

func ExampleBasic() {
	fmt.Println(ErrAccountExists)
	// Output:
	// The user account already exists
}
