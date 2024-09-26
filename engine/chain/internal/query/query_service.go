package query

import (
	"context"
	"errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"net"
)

type Service struct {
	logger       zerolog.Logger
	httpClient   *resty.Client
	url          string
	token        string
	deploymentID string
}

func (s *Service) Query(ctx context.Context, txID string) (*chainmodels.TXInfo, chainmodels.QueryTXOutcome, error) {
	req := s.httpClient.R()

	if s.token != "" {
		req.SetHeader("Content-Type", "application/json")
	}

	if s.deploymentID != "" {
		req.SetHeader("XDeployment-ID", s.deploymentID)
	}

	response, err := req.
		SetHeader("Content-Type", "application/json").
		Get(s.url + txID)

	if err != nil {
		var e net.Error
		if errors.As(err, &e) {
			return nil, chainmodels.QueryTxOutcomeFailed, spverrors.ErrARCUnreachable.Wrap(e)
		}
		return nil, chainmodels.QueryTxOutcomeFailed, spverrors.ErrInternal.Wrap(err)
	}

	response.Status()

	return nil, chainmodels.QueryTxOutcomeFailed, nil
}

func NewQueryService(logger zerolog.Logger, httpClient *resty.Client, url, token string) *Service {
	return &Service{
		logger:     logger,
		httpClient: httpClient,
		url:        url,
		token:      token,
	}
}
