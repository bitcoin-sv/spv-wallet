//go:build errorx
// +build errorx

package actions

import (
	"strings"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/spike/domain/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/spike/repos"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
)

type SomeEndpoint struct {
	service *transaction.Service
	logger  zerolog.Logger
}

func NewSomeEndpoint() *SomeEndpoint {
	repo := repos.NewRepo()

	service := transaction.NewService(repo)

	logger := zerolog.New(zerolog.NewConsoleWriter())

	return &SomeEndpoint{
		service: service,
		logger:  logger,
	}
}

func (e *SomeEndpoint) SomeSearch(c *gin.Context) Response {
	query := c.Query("query")

	r, err := e.service.Search(query)

	switch errorx.TypeSwitch(err, errorx.DataUnavailable, errorx.IllegalArgument, errorx.IllegalState, errorx.InternalError) {
	case nil:
		return Response{
			Status: 404,
			Body:   r,
		}
	case errorx.IllegalArgument:
		e.logger.Warn().Msgf("%+v", err)
		return Response{
			Status: 400,
			Body:   ErrorResponseOf(err),
		}
	case errorx.IllegalState:
		e.logger.Warn().Msgf("%+v", err)
		return Response{
			Status: 422,
			Body:   ErrorResponseOf(err),
		}
	case errorx.DataUnavailable:
		e.logger.Info().Msgf("%+v", err)
		return Response{
			Status: 404,
			Body:   ErrorResponseOf(err),
		}
	case errorx.InternalError:
		e.logger.Error().Err(err).Msgf("%+v", err)
		return Response{
			Status: 500,
			Body:   "internal server error",
		}
	default:
		e.logger.Error().Err(err).Msgf("unexpected error %+v", err)
		return Response{
			Status: 500,
			Body: ErrResponse{
				Code:    "unknown_error",
				Message: "Unexpected error",
			},
		}
	}
}

type Request struct {
	ID string
}

type Response struct {
	Status int
	Body   any
}

type ErrResponse struct {
	Code    string
	Message string
}

func ErrorResponseOf(err error) ErrResponse {
	errx := errorx.Cast(err)
	if errx == nil {
		panic("ErrorResponseOf must be used on errorx.Error")
	}
	return ErrResponse{
		Code:    shortName(errx.Type()),
		Message: errx.Message(),
	}
}

func shortName(errorType *errorx.Type) string {
	errorFullName := errorType.FullName()
	parent := errorType.Supertype().FullName()
	return strings.Replace(errorFullName, parent+".", "", 1)
}
