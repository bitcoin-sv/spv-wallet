package spverrors

import (
	"errors"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ErrorResponse is searching for error and setting it up in gin context
func ErrorResponse(c *gin.Context, err error, log *zerolog.Logger) {
	response, statusCode := getError(err, log)
	c.JSON(statusCode, response)
}

// AbortWithErrorResponse is searching for error and abort with error set
func AbortWithErrorResponse(c *gin.Context, err error, log *zerolog.Logger) {
	response, statusCode := getError(err, log)
	c.AbortWithStatusJSON(statusCode, response)
}

func getError(err error, log *zerolog.Logger) (models.ResponseError, int) {
	var errDetails models.SPVError
	ok := errors.As(err, &errDetails)
	if !ok {
		logError(log, err)
		return models.ResponseError{Code: models.UnknownErrorCode, Message: "Unable to get information about error"}, 500
	}

	return models.ResponseError{Code: errDetails.Code, Message: errDetails.Message}, errDetails.StatusCode
}

func logError(log *zerolog.Logger, err error) {
	if log != nil {
		log.Warn().Str("module", "spv-errors").Msgf("Unable to get information about error, details:  %s", err.Error())
	}
}
