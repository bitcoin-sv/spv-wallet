package spverrors

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Response struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

func ErrorResponse(c *gin.Context, err error, log *zerolog.Logger) {
	response, statusCode, found := getError(err)
	if found {
		log.Warn().Str("module", "spv-errors").Msgf("Unable to get information about error, details:  %s", err.Error())
	}
	c.JSON(statusCode, response)
}

func getError(err error) (Response, int, bool) {
	errDetails, ok := SPVErrorResponses[err]
	if !ok {
		return Response{Code: UnknownError, Message: "Unable to get information about error"}, 500, false
	}

	return Response{Code: errDetails.Code, Message: errDetails.Message}, errDetails.StatusCode, true
}
