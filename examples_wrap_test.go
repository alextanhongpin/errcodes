package errcodes_test

import (
	"fmt"

	"github.com/alextanhongpin/errcodes"
)

var ErrPayoutDeclined = errcodes.New(errcodes.Conflict, "payout_declined", "Payout cannot be processed. Please contact customer support for more information")

func ExampleWrap() {
	err := handlePayout()
	fmt.Println(errcodes.Sprint(err, false))
	fmt.Println()
	fmt.Println("Reversed:")
	fmt.Println(errcodes.Sprint(err, true))

	// Output:
	// Error: Payout cannot be processed. Please contact customer support for more information
	//     Origin is: account is actually frozen
	//         at errcodes_test.handlePayout (in examples_wrap_test.go:34)
	//     Ends here:
	//         at errcodes_test.ExampleWrap (in examples_wrap_test.go:12)
	//
	// Reversed:
	// Error: Payout cannot be processed. Please contact customer support for more information
	//     Ends here:
	//         at errcodes_test.ExampleWrap (in examples_wrap_test.go:12)
	//     Origin is: account is actually frozen
	//         at errcodes_test.handlePayout (in examples_wrap_test.go:34)
}

func handlePayout() error {
	err := errcodes.Wrap(ErrPayoutDeclined, "account is actually frozen")
	return err
}
