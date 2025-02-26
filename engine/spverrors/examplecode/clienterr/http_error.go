package clienterr

import (
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/errdef"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
)

func Response(c *gin.Context, err error, log *zerolog.Logger) {
	problem, logLevel := problemDetailsFromError(err)
	log.WithLevel(logLevel).Err(err).Msgf("Error HTTP response, returning %d: %s", problem.Status, problem.Detail)
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
	problem = errdef.ProblemDetails{
		Type:     "unknown_error",
		Title:    "Unknown error",
		Status:   500,
		Instance: "unknown_error",
	}
	return
}
