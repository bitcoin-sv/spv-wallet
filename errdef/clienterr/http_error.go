package clienterr

import (
	"errors"

	"github.com/bitcoin-sv/spv-wallet/errdef"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
)

// Response sends the error as a JSON response to the client.
func Response(c *gin.Context, err error, log *zerolog.Logger) {
	problem, logLevel := problemDetailsFromError(err)

	l := log.WithLevel(logLevel)
	if problem.Status >= 500 {
		l.Stack()
	}
	l.Err(err).Msgf("Error HTTP response, returning %d: %s", problem.Status, problem.Detail)

	c.JSON(problem.Status, problem)
}

func problemDetailsFromError(err error) (problem errdef.ProblemDetails, level zerolog.Level) {
	var ex *errorx.Error
	if errors.As(err, &ex) {
		if details, ok := ex.Property(propProblemDetails); ok {
			problem = details.(errdef.ProblemDetails)
			level = zerolog.InfoLevel
			return
		}

		// map internal error to problem details
		level = zerolog.WarnLevel
		problem.Type = "internal"
		problem.FromInternalError(ex)
		if errorx.HasTrait(ex, errdef.TraitUnsupported) {
			problem.Title = "Unsupported operation"
			problem.Status = 501
			return
		}
		if errorx.HasTrait(ex, errdef.TraitShouldNeverHappen) {
			problem.Detail = "This should never happen"
		}

		problem.Title = "Internal Server Error"
		problem.Status = 500
		return
	}

	level = zerolog.ErrorLevel
	problem.Title = "Unknown error"
	problem.Status = 500
	problem.Type = "internal"
	return
}
