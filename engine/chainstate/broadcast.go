package chainstate

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var (
	// broadcastSuccessErrors are a list of errors that are still considered a success
	broadcastSuccessErrors = []string{
		"already in the mempool", // {"error": "-27: Transaction already in the mempool"}
		"txn-already-know",       // { "error": "-26: 257: txn-already-known"}  // txn-already-know
		"txn-already-in-mempool", // txn-already-in-mempool
		"txn_already_known",      // TXN_ALREADY_KNOWN
		"txn_already_in_mempool", // TXN_ALREADY_IN_MEMPOOL
	}

	// broadcastQuestionableErrors are a list of errors that are not good broadcast responses,
	// but need to be checked differently
	broadcastQuestionableErrors = []string{
		"missing inputs", // Returned from mAPI for a valid tx that is on-chain
	}

	/*
		TXN_ALREADY_KNOWN (suppressed - returns as success: true)
		TXN_ALREADY_IN_MEMPOOL (suppressed - returns as success: true)
		TXN_MEMPOOL_CONFLICT
		NON_FINAL_POOL_FULL
		TOO_LONG_NON_FINAL_CHAIN
		BAD_TXNS_INPUTS_TOO_LARGE
		BAD_TXNS_INPUTS_SPENT
		NON_BIP68_FINAL
		TOO_LONG_VALIDATION_TIME
		BAD_TXNS_NONSTANDARD_INPUTS
		ABSURDLY_HIGH_FEE
		DUST
		TX_FEE_TOO_LOW
	*/
)

// broadcast will broadcast using a standard strategy
//
// NOTE: if successful (in-mempool), no error will be returned
func (c *Client) broadcast(ctx context.Context, id, hex string, format HexFormatFlag, timeout time.Duration, resultsChannel chan *BroadcastResult) {
	// Create a context (to cancel or timeout)
	ctxWithCancel, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var wg sync.WaitGroup

	for _, broadcastProvider := range createActiveProviders(c, id, hex, format) {
		wg.Add(1)
		go func(provider txBroadcastProvider) {
			defer wg.Done()
			resultsChannel <- broadcastToProvider(ctxWithCancel, ctx, provider, id, c, timeout)
		}(broadcastProvider)
	}

	wg.Wait()
	close(resultsChannel)
}

func createActiveProviders(c *Client, txID, txHex string, format HexFormatFlag) []txBroadcastProvider {
	providers := make([]txBroadcastProvider, 0, 1)

	switch c.ActiveProvider() {
	case ProviderMinercraft:
		if format != RawTx {
			panic("MAPI doesn't support other broadcast format than RawTx")
		}

		for _, miner := range c.options.config.minercraftConfig.broadcastMiners {
			if miner == nil {
				continue
			}

			pvdr := mapiBroadcastProvider{miner: miner, txID: txID, txHex: txHex}
			providers = append(providers, &pvdr)
		}
	case ProviderBroadcastClient:
		pvdr := broadcastClientProvider{txID: txID, txHex: txHex, format: format}
		providers = append(providers, &pvdr)
	default:
		c.options.logger.Warn().Msg("no active provider for broadcast")
	}

	return providers
}

func broadcastToProvider(ctx, fallbackCtx context.Context, provider txBroadcastProvider, txID string,
	c *Client, fallbackTimeout time.Duration,
) *BroadcastResult {
	failure := provider.broadcast(ctx, c)

	if failure != nil {
		checkMempool := containsAny(failure.Error.Error(), broadcastQuestionableErrors)

		if !checkMempool { // return original failure
			return &BroadcastResult{
				Provider: provider.getName(),
				Failure:  failure,
			}
		}

		// check in Mempool as fallback - if transaction is there -> GREAT SUCCESS
		if _, err := c.QueryTransaction(fallbackCtx, txID, requiredInMempool, fallbackTimeout); err != nil {
			return &BroadcastResult{
				Provider: provider.getName(),
				Failure: &BroadcastFailure{
					InvalidTx: failure.InvalidTx,
					Error:     fmt.Errorf("query tx failed: %w, initial error: %s", err, failure.Error.Error()),
				},
			}
		}
	}

	// successful broadcasted or found in mempool
	return &BroadcastResult{Provider: provider.getName()}
}
