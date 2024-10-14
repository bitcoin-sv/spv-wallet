package chainerrors

import (
	"encoding/json"
	"net/http"
	"strconv"

	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

// MapBHSErrorResponseToSpverror is a method that will check what kind of response came back from
// Block Header Service and map it to spverror and set it to context
func MapBHSErrorResponseToSpverror(res *resty.Response) models.SPVError {
	var responseErr chainmodels.BHSError

	if err := json.Unmarshal(res.Body(), &responseErr); err != nil {

		// TODO: remove after SPV-1106 is done
		if bodyStr := string(res.Body()); bodyStr != "" {
			// Try to unescape the string to remove any escaped characters like \" or \\n, etc.
			unescapedBodyStr, unescapeErr := strconv.Unquote(bodyStr)
			if unescapeErr != nil {
				// If unquoting fails, return the original string without modification
				unescapedBodyStr = bodyStr
			}

			return models.SPVError{
				Message:    unescapedBodyStr,
				StatusCode: http.StatusInternalServerError,
				Code:       spverrors.ErrInternal.Code,
			}.Wrap(err)
		}
		return ErrBHSParsingResponse.Wrap(err)
	}

	switch responseErr.Code {
	case "ErrInvalidBatchSize":
		return ErrInvalidBatchSize
	case "ErrMerkleRootNotFound":
		return ErrMerkleRootNotFound
	case "ErrMerkleRootNotInLC":
		return ErrMerkleRootNotInLongestChain
	case "error-unauthorized":
		return ErrBHSUnauthorized
	default:
		spvErr := models.SPVError{
			Message:    responseErr.Message,
			StatusCode: res.StatusCode(),
			Code:       responseErr.Code,
		}
		return spverrors.ErrInternal.Wrap(spvErr)
	}

}
