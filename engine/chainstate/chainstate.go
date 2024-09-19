/*
Package chainstate is the on-chain data service abstraction layer
*/
package chainstate

import (
	"context"
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

// BroadcastResult contains data about broadcasting to provider
type BroadcastResult struct {
	Provider string
	Failure  *BroadcastFailure
}

// BroadcastFailure contains data about broadcast failure
type BroadcastFailure struct {
	InvalidTx bool
	Error     error
}

// Broadcast will attempt to broadcast a transaction using the given providers
func (c *Client) Broadcast(ctx context.Context, id, txHex string, format HexFormatFlag, timeout time.Duration) *BroadcastResult {
	results := make(chan *BroadcastResult)
	go c.broadcast(ctx, id, txHex, format, timeout, results)

	failures := make([]*BroadcastResult, 0)

	for r := range results {
		if r.Failure != nil {
			failures = append(failures, r)
		} else {
			return r // one successful result is sufficient, and we consider the entire broadcast process complete. We disregard failures from other providers
		}
	}

	return groupBroadcastResults(failures)
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
