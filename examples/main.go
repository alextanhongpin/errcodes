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
	_               = errcodes.Load(errorBytes, toml.Unmarshal)
	ErrUserNotFound = errcodes.Of("user_not_found")
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
}
