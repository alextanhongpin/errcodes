package errcodes_test

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/errcodes"
)

type UserExistsError struct {
	error
	ID string
}

func NewUserExistsError(id string) error {
	return &UserExistsError{
		// Embed the error to allow comparison with errors.Is.
		error: ErrUserExists,

		// The additional properties we want to add.
		ID: id,
	}
}

func (e *UserExistsError) Error() string {
	return fmt.Sprintf("%s: %s", e.error.Error(), e.ID)
}

// Unwrap returns the wrapped error, hence allowing errors.Is to work.
func (e *UserExistsError) Unwrap() error {
	return e.error
}

func ExampleCustomData() {
	err := NewUserExistsError("user-42")
	fmt.Println(err)
	fmt.Println(errors.Is(err, ErrUserExists))

	var ec *errcodes.Error
	if errors.As(err, &ec) {
		fmt.Println(ec.String())
	}

	var userErr *UserExistsError
	if errors.As(err, &userErr) {
		fmt.Println(userErr.ID)
	}

	// Output:
	// The user account already exists: user-42
	// true
	// exists/user_exists: The user account already exists
	// user-42
}
