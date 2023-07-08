package errcodes_test

import (
	"fmt"

	"github.com/alextanhongpin/errcodes"
)

var ErrPayoutDeclined = errcodes.New(errcodes.Conflict, "payout_declined", "Payout cannot be processed. Please contact customer support for more information")

func ExampleWrap() {
	err := handlePayout()
	fmt.Println(errcodes.Sprint(err))
	fmt.Println()
	fmt.Println("Reversed:")
	fmt.Println(errcodes.Sprint(err, true))

	// Output:
	// Error: Payout cannot be processed. Please contact customer support for more information
	//     Caused by: account is actually frozen
	//         at errcodes_test.handlePayout (in examples_wrap_test.go:36)
	//         at errcodes_test.ExampleWrap (in examples_wrap_test.go:12)
	//     Ends here:
	//         at errcodes_test.ExampleWrap (in examples_wrap_test.go:13)
	//
	// Reversed:
	// Error: Payout cannot be processed. Please contact customer support for more information
	//     Ends here:
	//         at errcodes_test.ExampleWrap (in examples_wrap_test.go:16)
	//         at errcodes_test.ExampleWrap (in examples_wrap_test.go:12)
	//     Caused by: account is actually frozen
	//         at errcodes_test.handlePayout (in examples_wrap_test.go:36)
}

func handlePayout() error {
	err := errcodes.Wrap(ErrPayoutDeclined, "account is actually frozen")
	return err
}
