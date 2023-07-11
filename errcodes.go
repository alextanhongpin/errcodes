package errcodes

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
)

var ErrInvalidKind = errors.New("errcodes: invalid kind")

type Code string

type Kind string

const (
	Aborted            Kind = "aborted"
	BadRequest         Kind = "bad_request"
	Canceled           Kind = "cancelled"
	Conflict           Kind = "conflict"
	DataLoss           Kind = "data_loss"
	DeadlineExceeded   Kind = "deadline_exceeded"
	Exists             Kind = "exists"
	Forbidden          Kind = "forbidden"
	Internal           Kind = "internal"
	NotFound           Kind = "not_found"
	NotImplemented     Kind = "not_implemented"
	OutOfRange         Kind = "out_of_range"
	PreconditionFailed Kind = "precondition_failed"
	TooManyRequests    Kind = "too_many_requests"
	Unauthorized       Kind = "unauthorized"
	Unavailable        Kind = "unavailable"
	Unknown            Kind = "unknown"
)

func (c Kind) Valid() bool {
	switch c {
	case
		Aborted,
		BadRequest,
		Canceled,
		Conflict,
		DataLoss,
		DeadlineExceeded,
		Exists,
		Forbidden,
		Internal,
		NotFound,
		NotImplemented,
		OutOfRange,
		PreconditionFailed,
		TooManyRequests,
		Unauthorized,
		Unavailable,
		Unknown:
		return true
	default:
		return false
	}
}

type Error struct {
	kind    Kind
	code    Code
	message string
}

// New returns a new error with the given code, reason and description.
func New(kind Kind, code Code, message string) error {
	if !kind.Valid() {
		panic(ErrInvalidKind)
	}

	return &Error{
		kind:    kind,
		code:    code,
		message: message,
	}
}

// Error satisfies the error interface.
func (e *Error) Error() string {
	return e.message
}

func (e *Error) Kind() Kind {
	return e.kind
}

func (e *Error) Code() Code {
	return e.code
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) String() string {
	return fmt.Sprintf("%s/%s: %s", e.kind, e.code, e.message)
}

// Is checks if the error is of the same kind and same code.
func (e *Error) Is(err error) bool {
	var ec *Error
	if !errors.As(err, &ec) {
		return false
	}

	return e.kind == ec.kind && e.code == ec.code
}

var httpStatusByKind = map[Kind]int{
	Aborted:            http.StatusConflict,
	BadRequest:         http.StatusBadRequest,
	Canceled:           499, // client closed request.
	Conflict:           http.StatusConflict,
	DataLoss:           http.StatusInternalServerError,
	DeadlineExceeded:   http.StatusGatewayTimeout,
	Exists:             http.StatusConflict,
	Forbidden:          http.StatusForbidden,
	Internal:           http.StatusInternalServerError,
	NotFound:           http.StatusNotFound,
	NotImplemented:     http.StatusNotImplemented,
	OutOfRange:         http.StatusBadRequest,
	PreconditionFailed: http.StatusBadRequest,
	TooManyRequests:    http.StatusTooManyRequests,
	Unauthorized:       http.StatusUnauthorized,
	Unavailable:        http.StatusServiceUnavailable,
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

// https://chromium.googlesource.com/external/github.com/grpc/grpc/+/refs/tags/v1.21.4-pre1/doc/statuscodes.md
var grpcCodeByKind = map[Kind]codes.Code{
	Aborted:            codes.Aborted,
	BadRequest:         codes.InvalidArgument,
	Canceled:           codes.Canceled,
	Conflict:           codes.Aborted,
	DataLoss:           codes.DataLoss,
	DeadlineExceeded:   codes.DeadlineExceeded,
	Exists:             codes.AlreadyExists,
	Forbidden:          codes.PermissionDenied,
	Internal:           codes.Internal,
	NotFound:           codes.NotFound,
	NotImplemented:     codes.Unimplemented,
	OutOfRange:         codes.OutOfRange,
	PreconditionFailed: codes.FailedPrecondition,
	TooManyRequests:    codes.ResourceExhausted,
	Unauthorized:       codes.Unauthenticated,
	Unavailable:        codes.Unavailable,
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

var kindByGRPCCode = func() map[codes.Code]Kind {
	m := make(map[codes.Code]Kind)
	for k, v := range grpcCodeByKind {
		m[v] = k
	}
	return m
}()

// GRPCCodeToHTTP returns the HTTP code for the given grpc code.
func GRPCCodeToHTTP(code codes.Code) int {
	kind, ok := kindByGRPCCode[code]
	if !ok {
		return http.StatusInternalServerError
	}

	return HTTPStatusCode(kind)
}
