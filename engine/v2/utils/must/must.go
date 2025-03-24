// Package must provides a simple way to panic if an error is not nil.
// It should be used only in initialisation phase,
// in places where error is most probably problem with the code itself.
package must

import (
	"fmt"
)

// HaveNoError panics if the error is not nil.
func HaveNoError(err error) {
	if err != nil {
		panic(err)
	}
}

// HaveNoErrorf panics if the error is not nil, wrapping error with a message.
func HaveNoErrorf(err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if err != nil {
		panic(fmt.Errorf("%s, caused by: %w", msg, err))
	}
}

// BeTrue is a simple way to panic if the condition is false.
func BeTrue(condition bool, format string, args ...interface{}) {
	if !condition {
		panic(fmt.Sprintf(format, args...))
	}
}
