package chainstate

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/libsv/go-bc"
)

// query will try ALL providers in order and return the first "valid" response based on requirements
func (c *Client) query(ctx context.Context, id string, requiredIn RequiredIn,
	timeout time.Duration,
) (*TransactionInfo, error) {
	// Create a context (to cancel or timeout)
	ctxWithCancel, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	txInfo, err := queryBroadcastClient(
		ctxWithCancel, c, id,
	)
	if err != nil {
		return nil, err
	}
	if !checkRequirementArc(requiredIn, id, txInfo) {
		return nil, spverrors.ErrCouldNotFindTransaction
	}

	return txInfo, nil
}

// queryBroadcastClient will submit a query transaction request to a go-broadcast-client
func queryBroadcastClient(ctx context.Context, client ClientInterface, id string) (*TransactionInfo, error) {
	resp, err := client.BroadcastClient().QueryTransaction(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, spverrors.ErrBroadcastUnreachable
		}
		var arcError *broadcast.ArcError
		if errors.As(err, &arcError) {
			if arcError.IsRejectedTransaction() {
				return nil, spverrors.ErrBroadcastRejectedTransaction.Wrap(err)
			}
		}
		return nil, spverrors.ErrCouldNotFindTransaction.Wrap(err)
	}

	if resp == nil || !strings.EqualFold(resp.TxID, id) {
		return nil, spverrors.ErrTransactionIDMismatch
	}

	bump, err := bc.NewBUMPFromStr(resp.BaseTxResponse.MerklePath)
	if err != nil {
		return nil, spverrors.ErrBroadcastWrongBUMPResponse.Wrap(err)
	}

	return &TransactionInfo{
		BlockHash:   resp.BlockHash,
		BlockHeight: resp.BlockHeight,
		ID:          resp.TxID,
		Provider:    resp.Miner,
		TxStatus:    resp.TxStatus,
		BUMP:        bump,
	}, nil
}
