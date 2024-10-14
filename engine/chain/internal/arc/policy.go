package arc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// GetPolicy requests ARC server for the policy
func (s *Service) GetPolicy(ctx context.Context) (*chainmodels.Policy, error) {
	result := &chainmodels.Policy{}
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
		return nil, s.wrapARCError(spverrors.ErrARCUnauthorized, arcErr)
	case http.StatusNotFound:
		return nil, spverrors.ErrARCUnreachable
	default:
		return nil, s.wrapARCError(spverrors.ErrARCUnsupportedStatusCode, arcErr)
	}
}
