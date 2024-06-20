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
	response, statusCode, found := getError(err)
	logError(found, log, err)
	c.JSON(statusCode, response)
}

// AbortWithErrorResponse is searching for error and abort with error set
func AbortWithErrorResponse(c *gin.Context, err error, log *zerolog.Logger) {
	response, statusCode, found := getError(err)
	logError(found, log, err)
	c.AbortWithStatusJSON(statusCode, response)
}

func getError(err error) (responseError, int, bool) {
	var errDetails SPVError
	ok := errors.As(err, &errDetails)
	if !ok {
		return responseError{Code: "error-unknown", Message: "Unable to get information about error"}, 500, false
	}

	return responseError{Code: errDetails.Code, Message: errDetails.Message}, errDetails.StatusCode, true
}

func logError(found bool, log *zerolog.Logger, err error) {
	if !found && log != nil {
		log.Warn().Str("module", "spv-errors").Msgf("Unable to get information about error, details:  %s", err.Error())
	}
}
