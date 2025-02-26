package clienterr

import (
	"github.com/joomcode/errorx"
)

type ClientErrorDefinition struct {
	title    string
	typeName string
	httpCode int
}

func (c ClientErrorDefinition) Wrap(cause error, msg string, args ...any) *errorx.Error {
	return c.New().Wrap(cause, msg, args...).Err()
}

func (c ClientErrorDefinition) New() *Builder {
	b := &Builder{from: &c}
	b.problemDetails.Status = c.httpCode
	b.problemDetails.Title = c.title

	b.problemDetails.Type = c.typeName
	return b
}
