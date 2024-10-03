package query

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// Service for querying transactions.
type Service struct {
	logger     zerolog.Logger
	httpClient *resty.Client
	arcCfg     chainmodels.ARCConfig
}

// NewQueryService creates a new query service.
func NewQueryService(logger zerolog.Logger, httpClient *resty.Client, arcCfg chainmodels.ARCConfig) *Service {
	return &Service{
		logger:     logger,
		httpClient: httpClient,
		arcCfg:     arcCfg,
	}
}

// QueryTransaction a transaction.
func (s *Service) QueryTransaction(ctx context.Context, txID string) (*chainmodels.TXInfo, error) {
	result := &chainmodels.TXInfo{}
	errResult := &chainmodels.ArcError{}
	req := s.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(result).
		SetError(errResult)

	if s.arcCfg.Token != "" {
		req.SetHeader("Authorization", s.arcCfg.Token)
	}

	if s.arcCfg.DeploymentID != "" {
		req.SetHeader("XDeployment-ID", s.arcCfg.DeploymentID)
	}

	response, err := req.Get(fmt.Sprintf("%s/v1/tx/%s", s.arcCfg.URL, txID))

	if err != nil {
		var e net.Error
		if errors.As(err, &e) {
			return nil, spverrors.ErrARCUnreachable.Wrap(e)
		}
		return nil, spverrors.ErrInternal.Wrap(err)
	}

	switch response.StatusCode() {
	case http.StatusOK:
		return result, nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, s.withArcError(errResult, spverrors.ErrARCUnauthorized)
	case http.StatusNotFound:
		if !errResult.IsEmpty() {
			// ARC returns 404 when transaction is not found
			return nil, nil // By convention, nil is returned when transaction is not found
		}
		return nil, spverrors.ErrARCUnreachable
	case http.StatusConflict:
		return nil, s.withArcError(errResult, spverrors.ErrARCGenericError)
	}

	return nil, s.withArcError(errResult, spverrors.ErrARCUnsupportedStatusCode)
}

func (s *Service) withArcError(errResult *chainmodels.ArcError, baseError models.SPVError) error {
	if errResult == nil || errResult.IsEmpty() {
		return baseError
	}
	return baseError.Wrap(errResult)
}
