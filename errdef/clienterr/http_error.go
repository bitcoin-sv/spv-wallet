package clienterr

import (
	"errors"

	errdef2 "github.com/bitcoin-sv/spv-wallet/errdef"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
)

// Response sends the error as a JSON response to the client.
func Response(c *gin.Context, err error, log *zerolog.Logger) {
	problem, logLevel := problemDetailsFromError(err)
	log.WithLevel(logLevel).Err(err).Msgf("Error HTTP response, returning %d: %s", problem.Status, problem.Detail)
	c.JSON(problem.Status, problem)
}

func problemDetailsFromError(err error) (problem errdef2.ProblemDetails, level zerolog.Level) {
	var ex *errorx.Error
	if errors.As(err, &ex) {
		if details, ok := ex.Property(propProblemDetails); ok {
			problem = details.(errdef2.ProblemDetails)
			level = zerolog.InfoLevel
			return
		}

		// map internal error to problem details
		level = zerolog.WarnLevel
		problem.Type = "internal"
		problem.FromInternalError(ex)
		if errorx.HasTrait(ex, errdef2.TraitUnsupported) {
			problem.Title = "Unsupported operation"
			problem.Status = 501
			return
		}
		if errorx.HasTrait(ex, errdef2.TraitShouldNeverHappen) {
			problem.Detail = "This should never happen"
		}

		problem.Title = "Internal Server Error"
		problem.Status = 500
		return
	}

	level = zerolog.ErrorLevel
	problem = errdef2.ProblemDetails{
		Type:     "internal",
		Title:    "Unknown error",
		Status:   500,
		Instance: "",
	}
	return
}
