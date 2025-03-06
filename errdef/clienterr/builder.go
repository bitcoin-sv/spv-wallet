package clienterr

import (
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/errdef"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
)

var propProblemDetails = errorx.RegisterProperty("problem_details")
var clientError = errorx.NewNamespace("client").NewType("error")

type Builder struct {
	from           *ClientErrorDefinition
	problemDetails errdef.ProblemDetails
	cause          error
}

func (b *Builder) Wrap(cause error) *Builder {
	b.cause = cause
	b.problemDetails.
		FromInternalError(cause)

	return b
}

func (b *Builder) Detailed(errType string, detail string, args ...any) *Builder {
	b.problemDetails.Type = errType
	b.problemDetails.PushDetail(fmt.Sprintf(detail, args...))
	return b
}

func (b *Builder) Err() *errorx.Error {
	var err *errorx.Error
	if b.cause != nil {
		err = clientError.WrapWithNoMessage(b.cause)
	} else {
		err = clientError.NewWithNoMessage()
	}
	return err.WithProperty(propProblemDetails, b.problemDetails)
}

func (b *Builder) Response(c *gin.Context, log *zerolog.Logger) {
	Response(c, b.Err(), log)
}

func (b *Builder) WithInstance(parts ...any) *Builder {
	var sb strings.Builder
	for i, p := range parts {
		sb.WriteString(fmt.Sprintf("%v", p))
		if i < len(parts)-1 {
			sb.WriteString("/")
		}
	}
	b.problemDetails.Instance = sb.String()
	return b
}
