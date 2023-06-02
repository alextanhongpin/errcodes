package errcodes_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/alextanhongpin/errcodes"
	"google.golang.org/grpc/codes"
)

var ErrUserExists = errcodes.New(errcodes.Exists, "user_exists", "The user account already exists")

func TestErrcodes(t *testing.T) {
	var err error = ErrUserExists

	tests := make(map[string]bool)
	tests["kind match"] = ErrUserExists.Kind == errcodes.Exists
	tests["code match"] = ErrUserExists.Code == "user_exists"
	tests["errors.Is returns true for unwrapped error"] = errors.Is(err, ErrUserExists)
	tests["errors.Is returns true for wrapped error"] = errors.Is(fmt.Errorf("%w: John", err), ErrUserExists)
	tests["errors.Unwrap matches"] = errors.Unwrap(fmt.Errorf("%w: John", err)) == ErrUserExists

	var errC *errcodes.Error
	tests["errors.As returns true"] = errors.As(err, &errC)

	errD := errC.WithData(map[string]any{
		"email": "john.appleseed@mail.com",
	})

	tests["metadata is set correctly"] = errD.Data["email"] == "john.appleseed@mail.com"
	tests["correct http status code"] = errcodes.HTTPStatusCode(errC.Kind) == http.StatusConflict
	tests["correct grpc status code"] = errcodes.GRPCCode(errC.Kind) == codes.AlreadyExists

	for name, ok := range tests {
		name, ok := name, ok
		t.Run(name, func(t *testing.T) {
			if !ok {
				t.Fatal("want true, got false")
			}
		})
	}
}
