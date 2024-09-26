package merkleroots

import (
	"encoding/json"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

type bHSErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// mapBHSErrorResponseToSpverror is a method that will check what kind of response came back from
// Block Header Service and map it to spverror and set it to context
func mapBHSErrorResponseToSpverror(ctx *gin.Context, res *resty.Response, logger *zerolog.Logger) {
	var responseErr bHSErrorResponse

	err := json.Unmarshal(res.Body(), &responseErr)
	if err != nil {
		spverrors.ErrorResponse(ctx, spverrors.ErrBHSParsingResponse, logger)
		return
	}

	switch responseErr.Code {
	case "ErrInvalidBatchSize":
		err = spverrors.ErrBHSInvalidBatchSize
	case "ErrMerkleRootNotFound":
		err = spverrors.ErrBHSMerkleRootNotFound
	case "ErrMerkleRootNotInLC":
		err = spverrors.ErrBHSMerkleRootNotInLC
	default:
		spvErr := models.SPVError{
			Message:    responseErr.Message,
			StatusCode: res.StatusCode(),
			Code:       responseErr.Code,
		}
		err = spverrors.ErrInternal.Wrap(spvErr)
	}

	spverrors.ErrorResponse(ctx, err, logger)
}
