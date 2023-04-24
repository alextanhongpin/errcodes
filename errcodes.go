package errcodes

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
)

var (
	ErrInvalidCode     = errors.New("errcodes: invalid code")
	ErrDuplicateReason = errors.New("errcodes: duplicate reason")
)

var registry map[string]*Error

func init() {
	registry = make(map[string]*Error)
}

type Code string

const (
	AlreadyExists       Code = "already_exists"
	BadRequest          Code = "bad_request"
	Conflict            Code = "conflict"
	Forbidden           Code = "forbidden"
	Internal            Code = "internal"
	NotFound            Code = "not_found"
	PreconditionsFailed Code = "precondition_failed"
	Unauthorized        Code = "unauthorized"
	Unknown             Code = "unknown"
)

func (c Code) Valid() bool {
	switch c {
	case
		"already_exists",
		"bad_request",
		"conflict",
		"forbidden",
		"internal",
		"not_found",
		"precondition_failed",
		"unauthorized",
		"unknown":
		return true
	default:
		return false
	}
}

type Error struct {
	Code        Code   // Code, e.g. failed_preconditions
	Reason      string // Unique reason code, e.g. user_exists
	Description string // Human-readable error description.
	Metadata    map[string]any
}

func New(code Code, reason, description string) *Error {
	if !code.Valid() {
		panic(ErrInvalidCode)
	}

	return &Error{
		Code:        code,
		Reason:      reason,
		Description: description,
		Metadata:    make(map[string]any),
	}
}

func (e *Error) Error() string {
	return e.Description
}

func (e *Error) Is(err error) bool {
	var ec *Error
	if !errors.As(err, &ec) {
		return false
	}

	return e.Reason == ec.Reason
}

func (e *Error) WithMetadata(metadata map[string]any) *Error {
	ec := e.Clone()
	ec.Metadata = metadata
	return ec
}

func (e *Error) Clone() *Error {
	err := New(e.Code, e.Reason, e.Description)
	for k, v := range e.Metadata {
		err.Metadata[k] = v
	}

	return err
}

// Of returns a copy of the error from the registry with the given reason.
func Of(reason string) *Error {
	err := registry[reason]
	return err.Clone()
}

// Load loads the errors from an external file, which can be embedded using go
// embed.
func Load(data []byte, unmarshalFn func(data []byte, v any) error) int {
	descriptionByReasonByCode := make(map[Code]map[string]string)
	if err := unmarshalFn(data, &descriptionByReasonByCode); err != nil {
		panic(fmt.Errorf("errcodes: unmarshal error: %w", err))
	}

	for code, descriptionByReason := range descriptionByReasonByCode {
		if !code.Valid() {
			panic(ErrInvalidCode)
		}

		for reason, description := range descriptionByReason {
			_, ok := registry[reason]
			if ok {
				panic(fmt.Errorf("%w: %q is repeated", ErrDuplicateReason, reason))
			}
			registry[reason] = New(code, reason, description)
		}
	}

	return len(descriptionByReasonByCode)
}

func HTTPStatusCode(code Code) int {
	switch code {
	case "already_exists":
		return http.StatusConflict
	case "bad_request":
		return http.StatusBadRequest
	case "conflict":
		return http.StatusConflict
	case "forbidden":
		return http.StatusForbidden
	case "internal":
		return http.StatusInternalServerError
	case "not_found":
		return http.StatusNotFound
	case "precondition_failed":
		return http.StatusPreconditionFailed
	case "unauthorized":
		return http.StatusUnauthorized
	case "unknown":
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func GRPCCode(code Code) codes.Code {
	switch code {
	case "already_exists":
		return codes.AlreadyExists
	case "bad_request":
		return codes.InvalidArgument
	case "conflict":
		return codes.Aborted
	case "forbidden":
		return codes.PermissionDenied
	case "internal":
		return codes.Internal
	case "not_found":
		return codes.NotFound
	case "precondition_failed":
		return codes.FailedPrecondition
	case "unauthorized":
		return codes.Unauthenticated
	case "unknown":
		return codes.Unknown
	default:
		return codes.Internal
	}
}
