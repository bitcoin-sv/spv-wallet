package spverrors

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// SPVError is extended error which holds information about http status and code that should be returned
type SPVError struct {
	Code       string
	Message    string
	StatusCode int
}

// responseError is an error which will be returned in HTTP response
type responseError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

func (e SPVError) Error() string {
	return e.Message
}

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

func getError(err error, log *zerolog.Logger) (responseError, int) {
	var errDetails SPVError
	ok := errors.As(err, &errDetails)
	if !ok {
		logError(log, err)
		return responseError{Code: "error-unknown", Message: "Unable to get information about error"}, 500
	}

	return responseError{Code: errDetails.Code, Message: errDetails.Message}, errDetails.StatusCode
}

func logError(log *zerolog.Logger, err error) {
	if log != nil {
		log.Warn().Str("module", "spv-errors").Msgf("Unable to get information about error, details:  %s", err.Error())
	}
}
