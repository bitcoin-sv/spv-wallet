package spverrors

// Wrapf wraps an error with a formatted message
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return NewInternalError("Internal Server Error").WithDetailsf(format, args...).Wrap(err)
}

// Newf creates a new error with a message (which can be formatted)
func Newf(message string, args ...interface{}) error {
	return NewInternalError("Internal Server Error").WithDetailsf(message, args...)
}
