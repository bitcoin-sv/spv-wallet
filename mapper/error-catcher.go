package mapper

import "errors"

type ErrorCatcher struct {
	errors error
}

func (c *ErrorCatcher) OK() bool {
	return c.errors == nil
}

func (c *ErrorCatcher) Error() error {
	return c.errors
}

func (c *ErrorCatcher) Fail(err error) *ErrorCatcher {
	c.errors = errors.Join(c.errors, err)
	return c
}

func (c *ErrorCatcher) NotOK() bool {
	return !c.OK()
}

func NewErrorCatcher() *ErrorCatcher {
	return &ErrorCatcher{
		errors: nil,
	}
}

func Try[T, R any](catcher *ErrorCatcher, iteratee func(item T, index int) (R, error)) func(item T, index int) R {
	return func(item T, index int) R {
		res, err := iteratee(item, index)
		if err != nil {
			catcher.Fail(err)
		}
		return res
	}
}
