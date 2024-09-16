package chainstate

import (
	"context"
	"errors"
	"strings"
	"sync"
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

// fastestQuery will try ALL providers on once and return the fastest "valid" response based on requirements
func (c *Client) fastestQuery(ctx context.Context, id string, requiredIn RequiredIn,
	timeout time.Duration,
) *TransactionInfo {
	// The channel for the internal results
	resultsChannel := make(
		chan *TransactionInfo,
	) // All miners & WhatsOnChain

	// Create a context (to cancel or timeout)
	ctxWithCancel, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Loop each miner (break into a Go routine for each query)
	var wg sync.WaitGroup

	wg.Add(1)
	go func(ctx context.Context, client *Client, wg *sync.WaitGroup, id string, requiredIn RequiredIn) {
		defer wg.Done()
		if resp, err := queryBroadcastClient(
			ctx, client, id,
		); err == nil && checkRequirementArc(requiredIn, id, resp) {
			resultsChannel <- resp
		}
	}(ctxWithCancel, c, &wg, id, requiredIn)

	// Waiting for all requests to finish
	go func() {
		wg.Wait()
		close(resultsChannel)
	}()

	return <-resultsChannel
}

// queryBroadcastClient will submit a query transaction request to a go-broadcast-client
func queryBroadcastClient(ctx context.Context, client ClientInterface, id string) (*TransactionInfo, error) {
	resp, err := client.BroadcastClient().QueryTransaction(ctx, id)
	if err != nil {
		var arcError *broadcast.ArcError
		if errors.As(err, &arcError) {
			if arcError.IsRejectedTransaction() {
				return nil, spverrors.ErrBroadcastRejectedTransaction.Wrap(err)
			}
			return nil, spverrors.ErrCouldNotFindTransaction.Wrap(err)
		}
		return nil, spverrors.ErrBroadcastUnreachable.Wrap(err)
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
