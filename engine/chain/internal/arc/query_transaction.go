package arc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
)

// QueryTransaction a transaction.
func (s *Service) QueryTransaction(ctx context.Context, txID string) (*chainmodels.TXInfo, error) {
	result := &chainmodels.TXInfo{}
	arcErr := &chainmodels.ArcError{}
	req := s.prepareARCRequest(ctx).
		SetResult(result).
		SetError(arcErr)

	response, err := req.Get(fmt.Sprintf("%s/v1/tx/%s", s.arcCfg.URL, txID))

	if err != nil {
		return nil, s.wrapRequestError(err)
	}

	switch response.StatusCode() {
	case http.StatusOK:
		return result, nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, s.wrapARCError(chainerrors.ErrARCUnauthorized, arcErr)
	case http.StatusNotFound:
		if !arcErr.IsEmpty() {
			// ARC returns 404 when transaction is not found
			return nil, nil // By convention, nil is returned when transaction is not found
		}
		return nil, chainerrors.ErrARCUnreachable
	case http.StatusConflict:
		return nil, s.wrapARCError(chainerrors.ErrARCGenericError, arcErr)
	default:
		return nil, s.wrapARCError(chainerrors.ErrARCUnsupportedStatusCode, arcErr)
	}
}
