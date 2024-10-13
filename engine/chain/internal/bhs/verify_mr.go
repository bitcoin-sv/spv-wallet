package bhs

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/bitcoin-sv/go-paymail/spv"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// VerifyMerkleRoots verifies the merkle roots of the given transactions using BHS request
func (s *Service) VerifyMerkleRoots(ctx context.Context, merkleRoots []*spv.MerkleRootConfirmationRequestItem) (valid bool, err error) {
	confirmations, err := s.makeVerifyMerkleRootsRequest(ctx, merkleRoots)
	if err != nil {
		return false, err
	}

	switch confirmations.ConfirmationState {
	case chainmodels.MRConfirmed:
		return true, nil
	case chainmodels.MRUnableToVerify:
		s.logger.Warn().Msg("BHS is up but could not verify some merkle root(s). Defaulting to treat the provided merkle roots as valid.")
		return true, nil
	case chainmodels.MRInvalid:
		return false, nil
	default:
		return false, spverrors.ErrInternal.Wrap(spverrors.Newf("unexpected confirmation state"))
	}
}

func (s *Service) makeVerifyMerkleRootsRequest(ctx context.Context, merkleRoots []*spv.MerkleRootConfirmationRequestItem) (*chainmodels.MerkleRootsConfirmations, error) {
	if len(merkleRoots) == 0 {
		return nil, chainerrors.ErrBHSBadRequest.Wrap(spverrors.Newf("at least one merkleroot is required"))
	}
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
		return nil, chainerrors.ErrBHSBadRequest.Wrap(spverrors.Newf("BHS error message: %v", string(response.Body())))
	default:
		return nil, chainerrors.ErrBHSNoSuccessResponse.Wrap(spverrors.Newf("BHS response status code: %d", response.StatusCode()))
	}
}
