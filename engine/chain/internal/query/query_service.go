package query

import (
	"context"
	"errors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"net"
	"net/http"

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

// Query a transaction.
func (s *Service) Query(ctx context.Context, txID string) (*chainmodels.TXInfo, error) {
	req := s.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&chainmodels.TXInfo{}).
		SetError(&chainmodels.ArcError{})

	if s.token != "" {
		req.SetHeader("Content-Type", "application/json")
	}

	if s.deploymentID != "" {
		req.SetHeader("XDeployment-ID", s.deploymentID)
	}

	response, err := req.Get(s.url + txID)

	if err != nil {
		var e net.Error
		if errors.As(err, &e) {
			return nil, spverrors.ErrARCUnreachable.Wrap(e)
		}
		return nil, spverrors.ErrInternal.Wrap(err)
	}

	switch response.StatusCode() {
	case http.StatusOK:
		txInfo, ok := response.Result().(*chainmodels.TXInfo)
		if !ok {
			return nil, spverrors.ErrARCParseResponse
		}
		return txInfo, nil
	case http.StatusUnauthorized:
		return nil, s.withArcError(response, spverrors.ErrARCUnauthorized)
	case http.StatusNotFound:
		return nil, nil
	case http.StatusConflict:
		return nil, s.withArcError(response, spverrors.ErrARCGenericError)
	}

	return nil, s.withArcError(response, spverrors.ErrARCUnsupportedStatusCode)
}

func (s *Service) withArcError(response *resty.Response, baseError models.SPVError) error {
	if response == nil || response.Error() == nil {
		return spverrors.ErrInternal
	}

	//nolint:errorlint // We get a model returned from the response so errors. So errors.Is function is not relevant here
	arcErr, ok := response.Error().(*chainmodels.ArcError)
	if !ok {
		return baseError.Wrap(spverrors.ErrARCParseResponse)
	}

	return baseError.Wrap(arcErr)
}
