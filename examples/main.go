package main

import (
	_ "embed"
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/alextanhongpin/errcodes"
)

//go:embed errors.toml
var errorBytes []byte

var (
	_                   = errcodes.Load(errorBytes, toml.Unmarshal)
	ErrUserNotFound     = errcodes.For("user_not_found")
	ErrUserEmailInvalid = errcodes.Register("bad_request.user_email_invalid: Your email is invalid")
)

func main() {
	var err error = ErrUserNotFound
	fmt.Println(errors.Is(err, ErrUserNotFound))
	var errC *errcodes.Error
	if !errors.As(err, &errC) {
		panic("invalid error")
	}

	fmt.Println(errC)
	errM := errC.WithMetadata(map[string]any{
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
