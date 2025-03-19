package clienterr

import (
	"fmt"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/errdef"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
)

var propProblemDetails = errorx.RegisterProperty("problem_details")
var clientError = errorx.NewNamespace("client").NewType("error")

// Builder is a fluent API for building client (4xx) errors.
type Builder struct {
	from           *ClientErrorDefinition
	problemDetails errdef.ProblemDetails
	cause          error
}

// Wrap wraps the provided error as the cause of the client error.
func (b *Builder) Wrap(cause error) *Builder {
	b.cause = cause
	b.problemDetails.
		FromInternalError(cause)

	return b
}

// Detailed changes the error type and adds a detail message.
func (b *Builder) Detailed(errType string, detail string, args ...any) *Builder {
	b.problemDetails.Type = errType
	b.problemDetails.PushDetail(fmt.Sprintf(detail, args...))
	return b
}

// Err returns the client error as an errorx.Error (which also implements error interface).
func (b *Builder) Err() *errorx.Error {
	var err *errorx.Error
	if b.cause != nil {
		err = clientError.WrapWithNoMessage(b.cause)
	} else {
		err = clientError.NewWithNoMessage()
	}
	return err.WithProperty(propProblemDetails, b.problemDetails)
}

// Response sends the client error as a JSON response to the client.
func (b *Builder) Response(c *gin.Context, log *zerolog.Logger) {
	Response(c, b.Err(), log)
}

// WithInstance sets the instance field of the problem details.
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
