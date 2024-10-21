package spverrors

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const internalErrMessage = "Internal server error"
const genericErrorCode = "error"

/*
Error represents an error (by default internal, but can be used also for user's errors)
It implements ExtendedError interface so it can be returned in HTTP response
By default all internal errors have the same code, that's why errors.Is() won't distinguish them
If you want to make your error unique, you should set the code by .WithCode() method
Don't put sensitive information into the public message (by .WithHTTPStatus), it will be returned to the client
If you want to log some sensitive information, use .WithDetails() method, which accumulates the details.
The details will be logged but not returned to the client

NOTE: All errors are returned by-value, and also value receivers are used,
That's because we want to keep the errors immutable, so it can be shared safely

Example:

How to define distinct internal error:
var ErrSth = NewError("Something goes wrong").WithCode("error-sth-goes-wrong")

How to return that error:

	if err != nil {
		return ErrSth.WithDetails("With some deep debug info: %d", someVar).Wrap(err)
	}

How to check that error:

	if errors.Is(err, ErrSth) {
		// all errors of type Error with code "error-sth-goes-wrong" will be caught here
	}

How to create generic internal error:

	if sthGoesWrong {
		return NewError("Error message 1")
	}

	if sthElseGoesWrong {
		return NewError("Error message 2")
	}

# Note that for above two errors, errors.Is() will not distinguish them, they will be treated as the same error

How to create user's error:

var ErrUser = NewError("User did sth wrong").WithHTTPStatus(http.StatusBadRequest, "Bad request").WithCode("error-user-wrong-action")

For above example, http response will be 400 with string "Bad request" in the body (in json - see models.ResponseError)
*/
type Error struct {
	PublicMessage string
	HTTPStatus    int

	code    string
	details string
	cause   error
}

// NewError creates a new error with formatted details message and default HTTP status 500
func NewError(details string) Error {
	return Error{
		PublicMessage: internalErrMessage,
		HTTPStatus:    http.StatusInternalServerError,

		code:    genericErrorCode,
		details: details,
	}
}

// NewErrorf creates a new error with formatted details message and default HTTP status 500
func NewErrorf(format string, args ...any) Error {
	return NewError(fmt.Sprintf(format, args...))
}

// WithHTTPStatus sets the HTTP status and "public" message for the error
func (e Error) WithHTTPStatus(status int, message string) Error {
	e.HTTPStatus = status
	e.PublicMessage = message
	return e
}

// WithDetails adds details to the error message
func (e Error) WithDetails(details string, args ...any) Error {
	e.details += "; " + fmt.Sprintf(details, args...)
	return e
}

// Unwrap returns the "cause" error
func (e Error) Unwrap() error {
	return e.cause
}

// Wrap sets the "cause" error
func (e Error) Wrap(err error) Error {
	e.cause = err
	return e
}

// WithCode sets the code for the error
// By this code errors.Is will distinguish this error from others
func (e Error) WithCode(code string) Error {
	e.code = code
	return e
}

// GetMessage returns the public error message for http response
func (e Error) GetMessage() string {
	return e.PublicMessage
}

// GetStatusCode returns the HTTP status code for http response
func (e Error) GetStatusCode() int {
	return e.HTTPStatus
}

// Error returns the error message with details, satisfying the error interface
func (e Error) Error() string {
	return e.details
}

// GetCode returns the error code string
func (e Error) GetCode() string {
	return e.code
}

// StackTrace returns the error's stack trace if the cause implements it.
func (e Error) StackTrace() (trace errors.StackTrace) {
	err, ok := e.cause.(interface {
		StackTrace() errors.StackTrace
	})
	if !ok {
		return
	}

	trace = err.StackTrace()

	return
}

// Is checks if the target error is the same as the current error
func (e Error) Is(target error) bool {
	t, ok := target.(Error)
	if !ok {
		return false
	}
	return e.code == t.code
}
