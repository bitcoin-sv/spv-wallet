package dberrors

import (
	"errors"

	"github.com/joomcode/errorx"
	"gorm.io/gorm"
)

// QueryOrNotFoundError wraps the error with QueryFailed or NotFound error type depending on the GORM error type.
func QueryOrNotFoundError(err error, message string, args ...any) *errorx.Error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound.Wrap(err, message, args...)
	} else {
		return QueryFailed.Wrap(err, message, args...)
	}
}
