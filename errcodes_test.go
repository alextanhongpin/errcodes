package errcodes_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/alextanhongpin/errcodes"
	"google.golang.org/grpc/codes"
)

var ErrUserAlreadyExists = errcodes.New(errcodes.AlreadyExists, "user_already_exists", "The user account already exists")

func TestErrcodes(t *testing.T) {
	var err error = ErrUserAlreadyExists
	var errC *errcodes.Error

	testcases := make(map[string]bool)
	testcases["errors.Is returns true for unwrapped error"] = errors.Is(err, ErrUserAlreadyExists)
	testcases["errors.Is returns true for wrapped error"] = errors.Is(fmt.Errorf("%w: John", err), ErrUserAlreadyExists)
	testcases["errors.Unwrap matches"] = errors.Unwrap(fmt.Errorf("%w: John", err)) == ErrUserAlreadyExists
	testcases["errors.As returns true"] = errors.As(err, &errC)

	errM := errC.WithMetadata(map[string]any{
		"email": "john.appleseed@mail.com",
	})

	testcases["metadata is set correctly"] = errM.Metadata["email"] == "john.appleseed@mail.com"
	testcases["correct http status code"] = errcodes.HTTPStatusCode(errC.Code) == http.StatusConflict
	testcases["correct grpc status code"] = errcodes.GRPCCode(errC.Code) == codes.AlreadyExists

	for name, ok := range testcases {
		name, ok := name, ok
		t.Run(name, func(t *testing.T) {
			if !ok {
				t.Fatal("expected result to be true, got false")
			}
		})
	}
}

func TestRegister(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := errcodes.Register("not_found.user_not_found: User does not exists")
		if want, got := errcodes.NotFound, err.Code; want != got {
			t.Fatalf("want %v, got %v", want, got)
		}

		if want, got := "user_not_found", err.Reason; want != got {
			t.Fatalf("want %v, got %v", want, got)
		}

		if want, got := "User does not exists", err.Description; want != got {
			t.Fatalf("want %v, got %v", want, got)
		}
	})

	t.Run("failed", func(t *testing.T) {
		defer func() {
			if e := recover(); e != nil {
				err, ok := e.(error)
				if !ok {
					t.Fatalf("want error, got %v", e)
				}

				if !errors.Is(err, errcodes.ErrInvalidFormat) {
					t.Fatalf("want %v, got %v", errcodes.ErrInvalidFormat, err)
				}
			}
		}()

		_ = errcodes.Register("hello")
	})
}
