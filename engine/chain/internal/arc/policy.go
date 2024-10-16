package arc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
)

func (s *Service) getPolicy(ctx context.Context) (*Policy, error) {
	result := &Policy{}
	arcErr := &chainmodels.ArcError{}
	req := s.prepareARCRequest(ctx).
		SetResult(result).
		SetError(arcErr)

	response, err := req.Get(fmt.Sprintf("%s/v1/policy", s.arcCfg.URL))

	if err != nil {
		return nil, s.wrapRequestError(err)
	}

	switch response.StatusCode() {
	case http.StatusOK:
		return result, nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, s.wrapARCError(chainerrors.ErrARCUnauthorized, arcErr)
	case http.StatusNotFound:
		return nil, chainerrors.ErrARCUnreachable
	default:
		return nil, s.wrapARCError(chainerrors.ErrARCUnsupportedStatusCode, arcErr)
	}
}
