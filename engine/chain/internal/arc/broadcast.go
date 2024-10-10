package arc

import (
	"context"
	"errors"
	"fmt"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/ef"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/go-resty/resty/v2"
	"net/http"
)

// Custom ARC defined http status codes
const (
	StatusNotExtendedFormat             = 460
	StatusMalformedTx                   = 461
	StatusInvalidInputs                 = 462
	StatusMalformedTx2                  = 463 //for some reason ARC documentation has two status codes with the same description "Malformed transaction" (461, 463)
	StatusInvalidOutputs                = 464
	StatusFeeTooLow                     = 465
	StatusMinedAncestorsNotFoundInBEEF  = 467
	StatusInvalidBumpInBEEF             = 468
	StatusInvalidMerkleRoots            = 469
	StatusCumulativeFeeValidationFailed = 473
)

// Broadcast submits a transaction to the ARC server and returns the transaction info.
func (s *Service) Broadcast(ctx context.Context, tx *sdk.Transaction) (*chainmodels.TXInfo, error) {
	result := &chainmodels.TXInfo{}
	arcErr := &chainmodels.ArcError{}
	req := s.prepareARCRequest(ctx).
		SetResult(result).
		SetError(arcErr)

	s.setCallbackHeaders(req)

	txHex, err := s.prepareTxHex(ctx, tx)
	if err != nil {
		return nil, err
	}

	req.SetBody(requestBody{
		RawTx: txHex,
	})

	response, err := req.Post(fmt.Sprintf("%s/v1/tx", s.arcCfg.URL))

	if err != nil {
		return nil, s.wrapRequestError(err)
	}

	switch response.StatusCode() {
	case http.StatusOK:
		return result, nil
	case http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound:
		return nil, s.wrapARCError(spverrors.ErrARCUnauthorized, arcErr)
	case http.StatusConflict:
		return nil, s.wrapARCError(spverrors.ErrARCGenericError, arcErr)
	case http.StatusBadRequest,
		http.StatusUnprocessableEntity,
		StatusMalformedTx,
		StatusMalformedTx2,
		StatusInvalidInputs,
		StatusInvalidOutputs,
		StatusMinedAncestorsNotFoundInBEEF,
		StatusInvalidBumpInBEEF,
		StatusInvalidMerkleRoots:
		return nil, s.wrapARCError(spverrors.ErrARCUnprocessable, arcErr)
	case StatusNotExtendedFormat:
		return nil, s.wrapARCError(spverrors.ErrARCNotExtendedFormat, arcErr)
	case StatusFeeTooLow, StatusCumulativeFeeValidationFailed:
		return nil, s.wrapARCError(spverrors.ErrARCWrongFee, arcErr)
	default:
		return nil, s.wrapARCError(spverrors.ErrARCUnsupportedStatusCode, arcErr)
	}
}

type requestBody struct {
	// Even though the name suggests that it is a raw transaction,
	// it is actually a hex encoded transaction
	// and can be in Raw, Extended Format or BEEF format.
	RawTx string `json:"rawTx"`
}

func (s *Service) prepareTxHex(ctx context.Context, tx *sdk.Transaction) (string, error) {
	if s.arcCfg.TxsGetter == nil {
		s.logger.Warn().Msg("TransactionsGetter is not set, can't convert transaction to EFHex. Using raw transaction hex as a fallback.")
		return tx.Hex(), nil
	}
	converter := ef.NewConverter(s.arcCfg.TxsGetter)
	efHex, err := converter.Convert(ctx, tx)
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return "", spverrors.ErrEFConvertInterrupted.Wrap(err)
	}
	if err != nil {
		s.logger.Warn().Err(err).Msg("Failed to convert transaction to EFHex. Using raw transaction hex as a fallback.")
		return tx.Hex(), nil
	}
	return efHex, nil
}

func (s *Service) setCallbackHeaders(req *resty.Request) {
	cb := s.arcCfg.Callback
	if cb != nil && cb.URL != "" {
		req.SetHeader("X-CallbackUrl", cb.URL)

		if cb.Token != "" {
			req.SetHeader("X-CallbackToken", cb.Token)
		}
	}
}
