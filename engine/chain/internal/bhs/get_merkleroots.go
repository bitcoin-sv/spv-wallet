package bhs

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// GetMerkleRoots returns Merkle Roots from Block Header Service
func (s *Service) GetMerkleRoots(ctx context.Context, query url.Values) (*models.MerkleRootsBHSResponse, error) {
	bhsURL, err := s.createBHSURL("/chain/merkleroot")
	if err != nil {
		return nil, err
	}

	var response models.MerkleRootsBHSResponse
	req := s.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&response).
		SetQueryParamsFromValues(query)

	if s.bhsCfg.AuthToken != "" {
		req.SetAuthToken(s.bhsCfg.AuthToken)
	} else {
		s.logger.Warn().Msg("warning creating Block Headers Service url - auth token is not set. Some requests might not work")
	}

	res, err := req.Get(bhsURL.String())
	if err != nil {
		var e net.Error
		if errors.As(err, &e) {
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
	url, err := url.Parse(fmt.Sprintf("%s/api/v1%s", s.bhsCfg.URL, endpointPath))
	if err != nil {
		return nil, chainerrors.ErrBHSBadURL.Wrap(err)
	}

	return url, nil
}
