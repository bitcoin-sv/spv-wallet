/*
Package chainstate is the on-chain data service abstraction layer
*/
package chainstate

import (
	"context"
	"fmt"
	"time"
)

// Broadcast will attempt to broadcast a transaction using the given providers
func (c *Client) Broadcast(ctx context.Context, id, txHex string, timeout time.Duration) (string, error) {
	// Basic validation
	if len(id) < 50 {
		return "", ErrInvalidTransactionID
	} else if len(txHex) <= 0 { // todo: validate the tx hex
		return "", ErrInvalidTransactionHex
	}

	// Debug the id and hex
	c.DebugLog("tx_id: " + id)
	c.DebugLog("tx_hex: " + txHex)

	// Broadcast or die
	successCompleteCh := make(chan string)
	errorCh := make(chan string)

	go c.broadcast(ctx, id, txHex, timeout, successCompleteCh, errorCh)

	// wait for first success
	success := <-successCompleteCh
	if success != "" {
		return success, nil
	}

	// successCompleteCh closed without any values
	errorMessage := <-errorCh
	return ProviderAll, fmt.Errorf("broadcast failed, errors: %s", errorMessage)
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
