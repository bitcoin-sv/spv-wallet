package models

import (
	"github.com/pkg/errors"
)

type ExtendedError interface {
	error
	GetCode() string
	GetMessage() string
	GetStatusCode() int
	StackTrace() (trace errors.StackTrace)
}

// SPVError is extended error which holds information about http status and code that should be returned
type SPVError struct {
	Code       string
	Message    string
	StatusCode int
	cause      error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// ResponseError is an error which will be returned in HTTP response
type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

const UnknownErrorCode = "error-unknown"

// Error returns the error message string for SPVError, satisfying the error interface
func (e SPVError) Error() string {
	return e.Message
}

// GetCode returns the error code string for SPVError
func (e SPVError) GetCode() string {
	return e.Code
}

// GetMessage returns the error message string for SPVError
func (e SPVError) GetMessage() string {
	return e.Message
}

// GetStatusCode returns the error status code for SPVError
func (e SPVError) GetStatusCode() int {
	return e.StatusCode
}

// StackTrace returns the error's stack trace.
func (e SPVError) StackTrace() (trace errors.StackTrace) {
	err, ok := e.cause.(stackTracer)
	if !ok {
		return
	}

	trace = err.StackTrace()

	return
}

func (e SPVError) Unwrap() error {
	return e.cause
}

func (e SPVError) Wrap(err error) SPVError {
	e.cause = err
	return e
}

func (e SPVError) WithTrace(err error) SPVError {
	if st := stackTracer(nil); !errors.As(e.cause, &st) {
		return e.Wrap(errors.WithStack(err))
	}
	return e.Wrap(err)
}
