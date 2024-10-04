package bhs

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitcoin-sv/go-paymail/spv"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"net"
	"net/http"
)

/**
NOTE: When you switch from enabled use_beef to disabled and restart the server,
sometimes paymail capabilities are cached and VerifyMerkleRoots is called even though it shouldn't be.
*/

func (s *Service) VerifyMerkleRoots(ctx context.Context, merkleRoots []*spv.MerkleRootConfirmationRequestItem) (*chainmodels.MerkleRootsConfirmations, error) {
	result := &chainmodels.MerkleRootsConfirmations{}
	req := s.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(s.bhsCfg.AuthToken).
		SetBody(merkleRoots).
		SetResult(result)

	response, err := req.Post(fmt.Sprintf("%s/api/v1/chain/merkleroot/verify", s.bhsCfg.URL))

	if err != nil {
		var e net.Error
		if errors.As(err, &e) {
			return nil, chainerrors.ErrBHSUnreachable.Wrap(err)
		}
		return nil, spverrors.ErrInternal.Wrap(err)
	}

	switch response.StatusCode() {
	case http.StatusOK:
		return result, nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, chainerrors.ErrBHSUnauthorized
	case http.StatusNotFound:
		return nil, chainerrors.ErrBHSUnreachable
	case http.StatusBadRequest:
		// Note: in case of error, BHS returns a string (not json) with the error message
		// Most common error is "at least one merkleroot is required"
		return nil, chainerrors.ErrBHSBadRequest.Wrap(spverrors.Newf("BHS response status code: %v", string(response.Body())))
	default:
		return nil, chainerrors.ErrBHSNoSuccessResponse.Wrap(spverrors.Newf("BHS response status code: %d", response.StatusCode()))
	}
}
