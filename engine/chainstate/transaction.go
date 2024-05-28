package chainstate

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/libsv/go-bc"
)

// query will try ALL providers in order and return the first "valid" response based on requirements
func (c *Client) query(ctx context.Context, id string, requiredIn RequiredIn,
	timeout time.Duration,
) *TransactionInfo {
	// Create a context (to cancel or timeout)
	ctxWithCancel, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	switch c.ActiveProvider() {
	case ProviderBroadcastClient:
		resp, err := queryBroadcastClient(
			ctxWithCancel, c, id,
		)
		if err == nil && checkRequirementArc(requiredIn, id, resp) {
			return resp
		}
	default:
		c.options.logger.Warn().Msg("no active provider for query")
	}
	return nil // No transaction information found
}

// fastestQuery will try ALL providers on once and return the fastest "valid" response based on requirements
func (c *Client) fastestQuery(ctx context.Context, id string, requiredIn RequiredIn,
	timeout time.Duration,
) *TransactionInfo {
	// The channel for the internal results
	resultsChannel := make(
		chan *TransactionInfo,
		// len(c.options.config.mAPI.queryMiners)+2,
	) // All miners & WhatsOnChain

	// Create a context (to cancel or timeout)
	ctxWithCancel, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Loop each miner (break into a Go routine for each query)
	var wg sync.WaitGroup

	switch c.ActiveProvider() {
	case ProviderBroadcastClient:
		wg.Add(1)
		go func(ctx context.Context, client *Client, wg *sync.WaitGroup, id string, requiredIn RequiredIn) {
			defer wg.Done()
			if resp, err := queryBroadcastClient(
				ctx, client, id,
			); err == nil && checkRequirementArc(requiredIn, id, resp) {
				resultsChannel <- resp
			}
		}(ctxWithCancel, c, &wg, id, requiredIn)
	default:
		c.options.logger.Warn().Msg("no active provider for fastestQuery")
	}

	// Waiting for all requests to finish
	go func() {
		wg.Wait()
		close(resultsChannel)
	}()

	return <-resultsChannel
}

// queryBroadcastClient will submit a query transaction request to a go-broadcast-client
func queryBroadcastClient(ctx context.Context, client ClientInterface, id string) (*TransactionInfo, error) {
	client.DebugLog("executing request using " + ProviderBroadcastClient)
	if resp, failure := client.BroadcastClient().QueryTransaction(ctx, id); failure != nil {
		client.DebugLog("error executing request using " + ProviderBroadcastClient + " failed: " + failure.Error())
		return nil, errors.New(failure.Error())
	} else if resp != nil && strings.EqualFold(resp.TxID, id) {
		bump, err := bc.NewBUMPFromStr(resp.BaseTxResponse.MerklePath)
		if err != nil {
			return nil, err
		}
		return &TransactionInfo{
			BlockHash:   resp.BlockHash,
			BlockHeight: resp.BlockHeight,
			ID:          resp.TxID,
			Provider:    resp.Miner,
			TxStatus:    resp.TxStatus,
			BUMP:        bump,
			// it's not possible to get confirmations from broadcast client; zero would be treated as "not confirmed" that's why -1
			Confirmations: -1,
		}, nil
	}
	return nil, ErrTransactionIDMismatch
}
