package merkleroots

import (
	"encoding/json"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

type bhsErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// mapBHSErrorResponseToSpverror is a method that will check what kind of response came back from
// Block Header Service and map it to spverror and set it to context
func mapBHSErrorResponseToSpverror(res *resty.Response) error {
	var responseErr bhsErrorResponse

	if err := json.Unmarshal(res.Body(), &responseErr); err != nil {
		return ErrBHSParsingResponse.Wrap(err)
	}

	switch responseErr.Code {
	case "ErrInvalidBatchSize":
		return ErrInvalidBatchSize
	case "ErrMerkleRootNotFound":
		return ErrMerkleRootNotFound
	case "ErrMerkleRootNotInLC":
		return ErrMerkleRootNotInLongestChain
	default:
		spvErr := models.SPVError{
			Message:    responseErr.Message,
			StatusCode: res.StatusCode(),
			Code:       responseErr.Code,
		}
		return spverrors.ErrInternal.Wrap(spvErr)
	}

}
