package errcodes_test

import (
	"errors"
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
	testcases["errors.Is returns true"] = errors.Is(err, ErrUserAlreadyExists)
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
