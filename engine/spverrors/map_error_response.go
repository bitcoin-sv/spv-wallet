package spverrors

import (
	"errors"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ResponseMapper provides a fluent API for mapping errors to structured HTTP responses.
// It allows conditional error handling with predefined responses while ensuring proper logging.
type ResponseMapper interface {
	// If checks uses errors.Is to match the error to the provided one.
	If(errToMatch error) OnMatch
	// Else sets the default error to return if no match was found.
	Else(errToReturn models.ExtendedError)
	// Finalize logs the base error and sends the response to the client.
	Finalize()
}

// OnMatch provides a fluent API for defining the response to return when a match is found.
type OnMatch interface {
	Then(errToReturn models.ExtendedError) ResponseMapper
}

// MapResponse creates a new ResponseMapper instance.
func MapResponse(c *gin.Context, err error, log *zerolog.Logger) ResponseMapper {
	return &responseMapperBuilder{
		c:       c,
		baseErr: err,
		log:     log,
		final:   nil,
	}
}

type responseMapperBuilder struct {
	c       *gin.Context
	baseErr error
	log     *zerolog.Logger
	match   error
	final   models.ExtendedError
}

func (r *responseMapperBuilder) If(errToMatch error) OnMatch {
	if r.final != nil {
		return r
	}
	r.match = nil
	if errors.Is(r.baseErr, errToMatch) {
		r.match = errToMatch
	}
	return r
}

func (r *responseMapperBuilder) Then(errToReturn models.ExtendedError) ResponseMapper {
	if r.final != nil {
		return r
	}
	if r.match != nil {
		r.final = errToReturn
	}
	return r
}

func (r *responseMapperBuilder) Else(errToReturn models.ExtendedError) {
	if r.final == nil {
		r.final = errToReturn
	}
	r.Finalize()
}

func (r *responseMapperBuilder) Finalize() {
	r.logBaseError()

	if r.final != nil {
		r.c.JSON(r.final.GetStatusCode(), models.ResponseError{
			Code:    r.final.GetCode(),
			Message: r.final.GetMessage(),
		})
		return
	}

	r.c.JSON(500, models.ResponseError{
		Code:    models.UnknownErrorCode,
		Message: "Internal server error",
	})
}

func (r *responseMapperBuilder) logBaseError() {
	if r.log == nil {
		return
	}
	err := r.baseErr
	logLevel := zerolog.WarnLevel
	statusCode := 500

	var extendedErr models.ExtendedError
	if errors.As(err, &extendedErr) {
		statusCode = extendedErr.GetStatusCode()
		if extendedErr.GetStatusCode() >= 500 {
			logLevel = zerolog.ErrorLevel
		}
	}
	logInstance := r.log.WithLevel(logLevel).Str("module", "spv-errors")
	logInstance.Err(err).Msgf("Error HTTP response, returning %d", statusCode)
}
