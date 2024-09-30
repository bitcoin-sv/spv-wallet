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

// QueryTransaction a transaction.
func (s *Service) QueryTransaction(ctx context.Context, txID string) (*chainmodels.TXInfo, error) {
	if !s.validateTX(txID) {
		return nil, spverrors.ErrInvalidTransactionID
	}
	req := s.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&chainmodels.TXInfo{}).
		SetError(&chainmodels.ArcError{})

	if s.token != "" {
		req.SetHeader("Authorization", s.token)
	}

	if s.deploymentID != "" {
		req.SetHeader("XDeployment-ID", s.deploymentID)
	}

	response, err := req.Get(fmt.Sprintf("%s/v1/tx/%s", s.url, txID))

	if err != nil {
		var e net.Error
		if errors.As(err, &e) {
			return nil, spverrors.ErrARCUnreachable.Wrap(e)
		}
		return nil, spverrors.ErrInternal.Wrap(err)
	}

	switch response.StatusCode() {
	case http.StatusOK:
		return response.Result().(*chainmodels.TXInfo), nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, s.withArcError(response, spverrors.ErrARCUnauthorized)
	case http.StatusNotFound:
		if _, ok := s.asARCError(response); ok {
			// ARC returns 404 when transaction is not found
			return nil, nil // By convention, nil is returned when transaction is not found
		}
		return nil, spverrors.ErrARCUnreachable
	case http.StatusConflict:
		return nil, s.withArcError(response, spverrors.ErrARCGenericError)
	}

	return nil, s.withArcError(response, spverrors.ErrARCUnsupportedStatusCode)
}

func (s *Service) validateTX(txID string) bool {
	return len(txID) >= 50
}

func (s *Service) withArcError(response *resty.Response, baseError models.SPVError) error {
	arcErr, ok := s.asARCError(response)
	if !ok {
		return baseError.Wrap(spverrors.ErrARCParseResponse)
	}
	return baseError.Wrap(arcErr)
}

func (s *Service) asARCError(response *resty.Response) (*chainmodels.ArcError, bool) {
	if response == nil || response.Error() == nil {
		return nil, false
	}

	arcErr := response.Error().(*chainmodels.ArcError)
	if arcErr.IsEmpty() {
		return nil, false
	}
	return arcErr, true
}
