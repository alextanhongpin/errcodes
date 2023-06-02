package main

import (
	_ "embed"
	"errors"
	"fmt"

	"github.com/alextanhongpin/errcodes"
)

var (
	ErrUserNotFound     = errcodes.New(errcodes.NotFound, "user_not_found", "The user is not found")
	ErrUserEmailInvalid = errcodes.New(errcodes.BadRequest, "user_email_invalid", "Your email is invalid")
)

func main() {
	var err error = ErrUserNotFound
	fmt.Println(errors.Is(err, ErrUserNotFound))
	var errC *errcodes.Error
	if !errors.As(err, &errC) {
		panic("invalid error")
	}

	fmt.Println(errC)
	errM := errC.WithData(map[string]any{
		"Name": "john",
	})
	fmt.Printf("%#v\n", errM)
	fmt.Printf("%#v\n", ErrUserNotFound)
	fmt.Println(errM)
	fmt.Println(errM.String())

	err = ErrUserEmailInvalid
	fmt.Println(err)
	fmt.Printf("%#v\n", err)
}
