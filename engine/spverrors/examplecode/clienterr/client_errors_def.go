package clienterr

import (
	"github.com/joomcode/errorx"
)

type ClientErrorDefinition struct {
	title    string
	httpCode int
	errType  *errorx.Type
}

func (c ClientErrorDefinition) Wrap(cause error, msg string, args ...any) *errorx.Error {
	return c.New().Wrap(cause, msg, args...).Err()
}

func (c ClientErrorDefinition) New() *Builder {
	b := &Builder{from: &c}
	b.problemDetails.Status = c.httpCode
	b.problemDetails.Title = c.title

	b.problemDetails.Type = c.errType.FullName()
	return b
}

func RegisterSubtype(parent ClientErrorDefinition, typename string, title string) ClientErrorDefinition {
	return ClientErrorDefinition{
		title:    title,
		httpCode: parent.httpCode,
		errType:  parent.errType.NewSubtype(typename),
	}
}
