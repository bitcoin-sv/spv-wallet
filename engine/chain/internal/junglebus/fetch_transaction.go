package junglebus

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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

	switch response.StatusCode() {
	case http.StatusOK:
		tx, err := sdk.NewTransactionFromBytes(result.Transaction)
		if err != nil {
			return nil, chainerrors.ErrJunglebusParseTransaction.Wrap(err)
		}
		return tx, nil
	case http.StatusNotFound:
		textContent := string(response.Body())
		if strings.Contains(textContent, "tx-not-found") {
			return nil, chainerrors.ErrJunglebusTxNotFound
		}
		return nil, chainerrors.ErrJunglebusFailure.Wrap(spverrors.Newf("junglebus returned 404 with body %s", textContent))
	default:
		return nil, chainerrors.ErrJunglebusFailure.Wrap(spverrors.Newf("junglebus returned status code %d", response.StatusCode()))
	}
}
