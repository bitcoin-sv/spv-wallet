package clienterr

import (
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/errdef"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
)

var propProblemDetails = errorx.RegisterProperty("problem_details")

type Builder struct {
	from           *ClientErrorDefinition
	problemDetails errdef.ProblemDetails
	cause          error
	wrapMsg        string
	wrapArgs       []any
}

func (b *Builder) Wrap(cause error, msg string, args ...any) *Builder {
	b.cause = cause
	b.problemDetails.
		FromInternalError(cause).
		PushDetail(fmt.Sprintf(msg, args...))

	return b
}

func (b *Builder) Err() *errorx.Error {
	t := b.from.errType
	var err *errorx.Error
	if b.cause != nil {
		err = t.Wrap(b.cause, b.wrapMsg, b.wrapArgs...)
	} else {
		err = t.NewWithNoMessage()
	}
	return err.WithProperty(propProblemDetails, b.problemDetails)
}

func (b *Builder) Response(c *gin.Context, log *zerolog.Logger) {
	Response(c, b.Err(), log)
}

func (b *Builder) WithDetailf(detail string, args ...any) *Builder {
	b.problemDetails.PushDetail(fmt.Sprintf(detail, args...))
	return b
}

func (b *Builder) WithInstancef(instance string, args ...any) *Builder {
	b.problemDetails.Instance = instance
	return b
}
