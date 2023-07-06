package errcodes

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
)

var (
	ErrInvalidKind   = errors.New("errcodes: invalid kind")
	ErrDuplicateCode = errors.New("errcodes: duplicate code")
	ErrInvalidFormat = errors.New("errcodes: invalid format")
)

type Code string

type Kind string

const (
	BadRequest         Kind = "bad_request"
	Conflict           Kind = "conflict"
	Exists             Kind = "exists"
	Forbidden          Kind = "forbidden"
	Internal           Kind = "internal"
	NotFound           Kind = "not_found"
	PreconditionFailed Kind = "precondition_failed"
	Unauthorized       Kind = "unauthorized"
	Unknown            Kind = "unknown"
)

func (c Kind) Valid() bool {
	switch c {
	case
		BadRequest,
		Conflict,
		Exists,
		Forbidden,
		Internal,
		NotFound,
		PreconditionFailed,
		Unauthorized,
		Unknown:
		return true
	default:
		return false
	}
}

type Error struct {
	Kind    Kind
	Code    Code
	Message string
}

// New returns a new error with the given code, reason and description.
func New(kind Kind, code Code, message string) *Error {
	if !kind.Valid() {
		panic(ErrInvalidKind)
	}

	return &Error{
		Kind:    kind,
		Code:    code,
		Message: message,
	}
}

// Error satisfies the error interface.
func (e *Error) Error() string {
	return e.Message
}

func (e *Error) String() string {
	return fmt.Sprintf("[errcodes.%s] %s: %s", e.Kind, e.Code, e.Message)
}

// Is checks if the error is of the same kind and same code.
func (e *Error) Is(err error) bool {
	var ec *Error
	if !errors.As(err, &ec) {
		return false
	}

	return e.Kind == ec.Kind && e.Code == ec.Code
}

var httpStatusByKind = map[Kind]int{
	BadRequest:         http.StatusBadRequest,
	Conflict:           http.StatusConflict,
	Exists:             http.StatusConflict,
	Forbidden:          http.StatusForbidden,
	Internal:           http.StatusInternalServerError,
	NotFound:           http.StatusNotFound,
	PreconditionFailed: http.StatusPreconditionFailed,
	Unauthorized:       http.StatusUnauthorized,
	Unknown:            http.StatusInternalServerError,
}

// ̱HTTPStatusCode returns the HTTP status code for the given error code.
func HTTPStatusCode(kind Kind) int {
	status, ok := httpStatusByKind[kind]
	if !ok {
		return http.StatusInternalServerError
	}
	return status
}

var grpcCodeByKind = map[Kind]codes.Code{
	BadRequest:         codes.InvalidArgument,
	Conflict:           codes.Aborted,
	Exists:             codes.AlreadyExists,
	Forbidden:          codes.PermissionDenied,
	Internal:           codes.Internal,
	NotFound:           codes.NotFound,
	PreconditionFailed: codes.FailedPrecondition,
	Unauthorized:       codes.Unauthenticated,
	Unknown:            codes.Unknown,
}

// GRPCCode returns the gRPC code for the given error code.
func GRPCCode(kind Kind) codes.Code {
	code, ok := grpcCodeByKind[kind]
	if !ok {
		return codes.Internal
	}
	return code
}
