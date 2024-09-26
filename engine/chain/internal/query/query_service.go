package query

import (
	"context"
	"errors"
	"net"

	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// Service for querying transactions.
type Service struct {
	logger       zerolog.Logger
	httpClient   *resty.Client
	url          string
	token        string
	deploymentID string
}

// Query a transaction.
func (s *Service) Query(ctx context.Context, txID string) (*chainmodels.TXInfo, chainmodels.QueryTXOutcome, error) {
	req := s.httpClient.R().SetContext(ctx)

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

// NewQueryService creates a new query service.
func NewQueryService(logger zerolog.Logger, httpClient *resty.Client, url, token, deploymentID string) *Service {
	return &Service{
		logger:       logger,
		httpClient:   httpClient,
		url:          url,
		token:        token,
		deploymentID: deploymentID,
	}
}
