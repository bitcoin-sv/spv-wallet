package mapper

import "errors"

// ErrorCatcher is a struct that catches errors and stores them in errors field
type ErrorCatcher struct {
	errors error
}

// OK checks if ErrorCatcher didn't catch any errors
func (c *ErrorCatcher) OK() bool {
	return c.errors == nil
}

// Error returns joined errors cought by ErrorCatcher
func (c *ErrorCatcher) Error() error {
	return c.errors
}

// Catch catches error and joins it in the errors field
func (c *ErrorCatcher) Catch(err error) *ErrorCatcher {
	c.errors = errors.Join(c.errors, err)
	return c
}

// NotOK checks if ErrorCatcher caught any errors
func (c *ErrorCatcher) NotOK() bool {
	return !c.OK()
}

// NewErrorCatcher returns an instance of ErrorCatcher
func NewErrorCatcher() *ErrorCatcher {
	return &ErrorCatcher{
		errors: nil,
	}
}

// Try takes ErrorCatcher and a function mapper that can return an error.
// If the mapper function returns an error it is joined to the errors field of ErrorCatcher
func Try[T, R any](catcher *ErrorCatcher, iteratee func(item T, index int) (R, error)) func(item T, index int) R {
	return func(item T, index int) R {
		res, err := iteratee(item, index)
		if err != nil {
			catcher.Catch(err)
		}
		return res
	}
}
