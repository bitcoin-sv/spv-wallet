package spverrors

import (
	"fmt"

	"github.com/pkg/errors"
)

// Wrapf wraps an error with a formatted message
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if len(args) > 0 {
		return errors.Wrapf(err, format+" caused by", args...)
	}
	return errors.Wrap(err, format+" caused by")
}

// Newf creates a new error with a message (which can be formatted)
func Newf(message string, args ...interface{}) error {
	if len(args) > 0 {
		return fmt.Errorf(message, args...)
	}
	return errors.New(message)
}
