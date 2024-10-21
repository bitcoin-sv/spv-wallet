package spverrors

// Wrapf wraps an error with a formatted message
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return NewErrorf(format, args...).Wrap(err)
}

// Newf creates a new error with a message (which can be formatted)
func Newf(message string, args ...interface{}) error {
	return NewErrorf(message, args...)
}

// Of creates spverrors.Error based on the error message
func Of(err error) Error {
	return NewError(err.Error())
}
