/*
Package chainstate is the on-chain data service abstraction layer
*/
package chainstate

import (
	"context"
	"fmt"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// HexFormatFlag transaction hex format
type HexFormatFlag byte

const (
	// RawTx is the raw transaction format
	RawTx HexFormatFlag = 1 << iota // 1
	// Ef is the Extended transaction format
	Ef
)

// Contains checks if the flag contains specific bytes
func (flag HexFormatFlag) Contains(other HexFormatFlag) bool {
	return (flag & other) == other
}

// SupportedBroadcastFormats returns supported formats based on active providers
func (c *Client) SupportedBroadcastFormats() HexFormatFlag {
	return RawTx | Ef
}

// BroadcastFailure contains data about broadcast failure
type BroadcastFailure struct {
	InvalidTx bool
	Error     error
}

// Broadcast will attempt to broadcast a transaction using the given providers
func (c *Client) Broadcast(ctx context.Context, txID, txHex string, format HexFormatFlag, timeout time.Duration) *BroadcastFailure {
	ctxWithCancel, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	failure := c.submitTransaction(ctxWithCancel, txID, txHex, format)

	if failure != nil {
		checkMempool := containsAny(failure.Error.Error(), broadcastQuestionableErrors)

		if !checkMempool { // return original failure
			return failure
		}

		// check in Mempool as fallback - if transaction is there -> GREAT SUCCESS
		if _, err := c.QueryTransaction(ctx, txID, requiredInMempool, timeout); err != nil {
			return &BroadcastFailure{
				InvalidTx: failure.InvalidTx,
				Error:     fmt.Errorf("query tx failed: %w, initial error: %s", err, failure.Error.Error()),
			}
		}
	}

	// successful broadcasted or found in mempool
	return nil
}

// QueryTransaction will get the transaction info from all providers returning the "first" valid result
//
// Note: this is slow, but follows a specific order: ARC -> WhatsOnChain
func (c *Client) QueryTransaction(
	ctx context.Context, id string, requiredIn RequiredIn, timeout time.Duration,
) (transaction *TransactionInfo, err error) {
	if c.options.metrics != nil {
		end := c.options.metrics.TrackQueryTransaction()
		defer func() {
			success := err == nil
			end(success)
		}()
	}

	// Basic validation
	if len(id) < 50 {
		return nil, spverrors.ErrInvalidTransactionID
	} else if !c.validRequirement(requiredIn) {
		return nil, spverrors.ErrInvalidRequirements
	}

	// Try all providers and return the "first" valid response
	info := c.query(ctx, id, requiredIn, timeout)
	if info == nil {
		return nil, spverrors.ErrCouldNotFindTransaction
	}
	return info, nil
}

// QueryTransactionFastest will get the transaction info from ALL provider(s) returning the "fastest" valid result
//
// Note: this is fast but could abuse each provider based on how excessive this method is used
func (c *Client) QueryTransactionFastest(
	ctx context.Context, id string, requiredIn RequiredIn, timeout time.Duration,
) (*TransactionInfo, error) {
	// Basic validation
	if len(id) < 50 {
		return nil, spverrors.ErrInvalidTransactionID
	} else if !c.validRequirement(requiredIn) {
		return nil, spverrors.ErrInvalidRequirements
	}

	// Try all providers and return the "fastest" valid response
	info := c.fastestQuery(ctx, id, requiredIn, timeout)
	if info == nil {
		return nil, spverrors.ErrCouldNotFindTransaction
	}
	return info, nil
}
