package spverrors

import (
	"errors"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ErrorResponse is searching for error and setting it up in gin context
func ErrorResponse(c *gin.Context, err error, log *zerolog.Logger) {
	response, statusCode := mapAndLog(err, log)
	c.JSON(statusCode, response)
}

// AbortWithErrorResponse is searching for error and abort with error set
func AbortWithErrorResponse(c *gin.Context, err error, log *zerolog.Logger) {
	response, statusCode := mapAndLog(err, log)
	c.AbortWithStatusJSON(statusCode, response)
}

func mapAndLog(err error, log *zerolog.Logger) (model models.ResponseError, statusCode int) {
	model.Code = models.UnknownErrorCode
	model.Message = "Internal server error"
	statusCode = 500

	var extendedErr models.ExtendedError
	if errors.As(err, &extendedErr) {
		model.Code = extendedErr.GetCode()
		model.Message = extendedErr.GetMessage()
		statusCode = extendedErr.GetStatusCode()
	}

	if log != nil {
		log.Warn().Str("module", "spv-errors").Err(err).Msgf("non-ExtendedError in HTTP response, returning %d", statusCode)
	}
	return
}
