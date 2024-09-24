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

	logLevel := zerolog.WarnLevel
	exposedInternalError := false
	var extendedErr models.ExtendedError
	if errors.As(err, &extendedErr) {
		model.Code = extendedErr.GetCode()
		model.Message = extendedErr.GetMessage()
		statusCode = extendedErr.GetStatusCode()
		if statusCode >= 500 {
			logLevel = zerolog.ErrorLevel
		}
	} else {
		// we should wrap all internal errors into SPVError (with proper code, message and status code)
		// if you find out that some endpoint produces this warning, feel free to fix it
		exposedInternalError = true
	}

	if log != nil {
		logInstance := log.WithLevel(logLevel).Str("module", "spv-errors")
		if exposedInternalError {
			logInstance.Str("warning", "internal error returned as HTTP response")
		}
		logInstance.Err(err).Msgf("Error HTTP response, returning %d", statusCode)
	}
	return
}
