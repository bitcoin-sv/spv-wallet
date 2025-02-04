package lox

import "errors"

// ErrorCollector is a struct that catches errors and stores them in errors field
type ErrorCollector struct {
	errors error
}

// OK checks if ErrorCollector didn't catch any errors
func (c *ErrorCollector) OK() bool {
	return c.errors == nil
}

// Error returns joined errors caught by ErrorCollector
func (c *ErrorCollector) Error() error {
	return c.errors
}

// Collect collects error and joins it in the errors field
func (c *ErrorCollector) Collect(err error) *ErrorCollector {
	c.errors = errors.Join(c.errors, err)
	return c
}

// NewErrorCollector returns an instance of ErrorCollector
func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		errors: nil,
	}
}
