package arc

import (
	"context"
	"errors"
	"net"

	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

func (s *Service) prepareARCRequest(ctx context.Context) *resty.Request {
	req := s.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json")

	if s.arcCfg.Token != "" {
		req.SetHeader("Authorization", s.arcCfg.Token)
	}

	if s.arcCfg.DeploymentID != "" {
		req.SetHeader("XDeployment-ID", s.arcCfg.DeploymentID)
	}

	return req
}

func (s *Service) wrapRequestError(err error) error {
	var e net.Error
	if errors.As(err, &e) {
		return spverrors.ErrARCUnreachable.Wrap(e)
	}
	return spverrors.ErrInternal.Wrap(err)
}

func (s *Service) wrapARCError(baseError models.SPVError, errResult *chainmodels.ArcError) error {
	if errResult == nil || errResult.IsEmpty() {
		return baseError
	}
	return baseError.Wrap(errResult)
}
