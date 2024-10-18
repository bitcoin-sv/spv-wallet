package spverrors

import (
	"fmt"

	"github.com/pkg/errors"
)

const internalErrCode = "internal-error"

/*
InternalError represents an internal error
It implements ExtendedError interface so it can be returned in HTTP response
By default all internal errors have the same code, that's why errors.Is() won't distinguish them
If you want to make your error unique, you should set the code by .WithCode() method
Don't put sensitive information into the public message, it will be returned to the client
For more detailed/debug information use .WithDetails() method
The details will be logged but not returned to the client
Example:

How to define distinct internal error:
var ErrSth = NewInternalError("Something goes wrong").WithCode("error-sth-goes-wrong")

How to return that error:

	if err != nil {
		return ErrSth.WithDetailsf("Some deep debug info %s", txID).Wrap(err)
	}

How to check that error:

	if errors.Is(err, ErrSth) {
		// all errors of type InternalError with code "error-sth-goes-wrong" will be caught here
	}

How to create generic internal error:

	if sthGoesWrong {
		return NewInternalError("Error message 1")
	}

	if sthElseGoesWrong {
		return NewInternalError("Error message 2")
	}

Note that for above two errors, errors.Is() will not distinguish them, they will be treated as the same error
*/
type InternalError struct {
	PublicMessage string
	Code          string

	details string
	cause   error
}

func NewInternalError(publicMessage string) InternalError {
	return InternalError{
		PublicMessage: publicMessage,
		Code:          internalErrCode,
	}
}

func (e InternalError) Unwrap() error {
	return e.cause
}

func (e InternalError) Wrap(err error) InternalError {
	e.cause = err
	return e
}

func (e InternalError) WithDetails(details string) InternalError {
	e.details = details
	return e
}

func (e InternalError) WithDetailsf(format string, args ...any) InternalError {
	e.details = fmt.Sprintf(format, args...)
	return e
}

func (e InternalError) WithCode(code string) InternalError {
	e.Code = code
	return e
}

func (e InternalError) GetMessage() string {
	return e.PublicMessage
}

func (e InternalError) GetStatusCode() int {
	return 500
}

func (e InternalError) Error() string {
	return e.PublicMessage
}

func (e InternalError) GetCode() string {
	return e.Code
}

func (e InternalError) StackTrace() (trace errors.StackTrace) {
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
func (e InternalError) Is(target error) bool {
	//nolint:errorlint //errors.Is/As would check also the wrapped error but here only the current one should be concerned
	t, ok := target.(InternalError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}
