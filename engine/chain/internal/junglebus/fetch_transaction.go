package junglebus

import (
	"context"
	"fmt"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// FetchTransaction fetches transaction from junglebus
// This method should be used only when there is no other way to get transaction data because it uses external service
func (s *Service) FetchTransaction(ctx context.Context, txID string) (*sdk.Transaction, error) {
	result := &TransactionResponse{}
	req := s.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(result)

	response, err := req.Get(fmt.Sprintf("https://junglebus.gorillapool.io/v1/transaction/get/%s", txID))

	if err != nil {
		return nil, spverrors.ErrInternal.Wrap(err)
	}

	if response.StatusCode() != 200 {
		return nil, chainerrors.ErrJunglebusFailure.Wrap(spverrors.Newf("junglebus returned status code %d", response.StatusCode()))
	}

	tx, err := sdk.NewTransactionFromBytes(result.Transaction)
	if err != nil {
		return nil, chainerrors.ErrJunglebusParseTransaction.Wrap(err)
	}

	return tx, nil
}
