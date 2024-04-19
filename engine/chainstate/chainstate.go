/*
Package chainstate is the on-chain data service abstraction layer
*/
package chainstate

import (
	"context"
	"time"
)

// HexFormatFlag transaction hex format
type HexFormatFlag byte

const (
	RawTx HexFormatFlag = 1 << iota // 1
	Ef
)

// Contains checks if the flag contains specific bytes
func (flag HexFormatFlag) Contains(other HexFormatFlag) bool {
	return (flag & other) == other
}

// SupportedBroadcastFormats retuns supported formats based on active providers
func (c *Client) SupportedBroadcastFormats() HexFormatFlag {
	switch c.ActiveProvider() {
	case ProviderMinercraft:
		return RawTx

	case ProviderBroadcastClient:
		return RawTx | Ef

	default:
		return RawTx
	}
}

type BroadcastResult struct {
	Provider string
	Failure  *BroadcastFailure
}
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
			return r // wait for first success
		}
	}

	return groupBroadcastResults(failures)
}

// QueryTransaction will get the transaction info from all providers returning the "first" valid result
//
// Note: this is slow, but follows a specific order: mAPI -> WhatsOnChain
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
		return nil, ErrInvalidTransactionID
	} else if !c.validRequirement(requiredIn) {
		return nil, ErrInvalidRequirements
	}

	// Try all providers and return the "first" valid response
	info := c.query(ctx, id, requiredIn, timeout)
	if info == nil {
		return nil, ErrTransactionNotFound
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
		return nil, ErrInvalidTransactionID
	} else if !c.validRequirement(requiredIn) {
		return nil, ErrInvalidRequirements
	}

	// Try all providers and return the "fastest" valid response
	info := c.fastestQuery(ctx, id, requiredIn, timeout)
	if info == nil {
		return nil, ErrTransactionNotFound
	}
	return info, nil
}
