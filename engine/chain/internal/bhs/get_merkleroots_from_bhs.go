package bhs

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"syscall"

	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// GetMerkleRootsFromBHS returns Merkle Roots from Block Header Service
func (s *Service) GetMerkleRootsFromBHS(ctx context.Context, query url.Values) (*models.MerkleRootsBHSResponse, error) {
	bhsURL, err := s.createBHSURL("/chain/merkleroot")
	if err != nil {
		return nil, err
	}

	req := s.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		EnableTrace()

	if s.bhsCfg.AuthToken != "" {
		req.SetAuthToken(s.bhsCfg.AuthToken)
	} else {
		s.logger.Warn().Msg("warning creating Block Headers Service url - auth token is not set. Some requests might not work")
	}

	var response models.MerkleRootsBHSResponse
	res, err := req.
		SetResult(&response).
		SetQueryParamsFromValues(query).
		Get(bhsURL.String())
	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) {
			return nil, chainerrors.ErrBHSUnreachable.Wrap(err)
		}
		return nil, spverrors.ErrInternal.Wrap(err)
	}
	if !res.IsSuccess() {
		return nil, mapBHSErrorResponseToSpverror(res)
	}

	return &response, nil
}

// createBHSURL parses Block Header Url from configuration and constructs a valid
// endpoint with provided endpointPath variable
func (s *Service) createBHSURL(endpointPath string) (*url.URL, error) {
	if s.bhsCfg.URL == "" {
		s.logger.Error().Msgf("create Block Header Service URL - url not configured")
	}

	url, err := url.Parse(fmt.Sprintf("%s/api/v1%s", s.bhsCfg.URL, endpointPath))
	if err != nil {
		return nil, chainerrors.ErrBHSBadURL.Wrap(err)
	}

	return url, nil
}
