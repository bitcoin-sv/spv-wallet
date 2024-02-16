package chainstate

import (
	"context"
	"fmt"
	"strings"
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
// NOTE: function register the fastest successful broadcast into 'completeChannel' so client doesn't need to wait for other providers
func (c *Client) broadcast(ctx context.Context, id, hex string, timeout time.Duration, completeChannel, errorChannel chan string) {
	// Create a context (to cancel or timeout)
	ctxWithCancel, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var wg sync.WaitGroup

	resultsChannel := make(chan broadcastResult)
	status := newBroadcastStatus(completeChannel)

	for _, broadcastProvider := range createActiveProviders(c, id, hex) {
		wg.Add(1)
		go func(provider txBroadcastProvider) {
			defer wg.Done()
			broadcastToProvider(ctxWithCancel, ctx, provider, id, c, timeout,
				resultsChannel, status)
		}(broadcastProvider)
	}

	go func() {
		wg.Wait()
		close(resultsChannel)
		status.dispose()
	}()

	var errorMessages []string
	for result := range resultsChannel {
		if result.isError {
			debugLog(c, id, fmt.Sprintf("broadcast error: %s from provider %s", result.err, result.provider))
			errorMessages = append(errorMessages, result.provider+": "+result.err.Error())
		} else {
			debugLog(c, id, fmt.Sprintf("successful broadcast to %s", result.provider))
		}
	}

	if !status.success && len(errorMessages) > 0 {
		errorChannel <- strings.Join(errorMessages, ", ")
	}
}

func createActiveProviders(c *Client, txID, txHex string) []txBroadcastProvider {
	providers := make([]txBroadcastProvider, 0, 10)

	switch c.ActiveProvider() {
	case ProviderMinercraft:
		for _, miner := range c.options.config.minercraftConfig.broadcastMiners {
			if miner == nil {
				continue
			}

			pvdr := mapiBroadcastProvider{miner: miner, txID: txID, txHex: txHex}
			providers = append(providers, &pvdr)
		}
	case ProviderBroadcastClient:
		pvdr := broadcastClientProvider{txID: txID, txHex: txHex}
		providers = append(providers, &pvdr)
	default:
		c.options.logger.Warn().Msg("no active provider for broadcast")
	}

	return providers
}

func broadcastToProvider(ctx, fallbackCtx context.Context, provider txBroadcastProvider, txID string,
	c *Client, fallbackTimeout time.Duration,
	resultsChannel chan broadcastResult, status *broadcastStatus,
) {
	bErr := provider.broadcast(ctx, c)

	if bErr != nil {
		// check in Mempool as fallback - if transaction is there -> GREAT SUCCESS
		// Check error response for "questionable errors"/(TX FAILURE)
		if doesErrorContain(bErr.Error(), broadcastQuestionableErrors) {
			bErr = checkInMempool(fallbackCtx, c, txID, bErr.Error(), fallbackTimeout)
		}

		if bErr != nil {
			resultsChannel <- newErrorResult(bErr, provider.getName())
		}
	}

	// successful broadcast or found in mempool
	if bErr == nil {
		status.tryCompleteWithSuccess(provider.getName())
		resultsChannel <- newSuccessResult(provider.getName())
	}
}

// checkInMempool is a quick check to see if the tx is in mempool (or on-chain)
func checkInMempool(ctx context.Context, client ClientInterface, id, initErrMsg string, timeout time.Duration) error {
	if _, err := client.QueryTransaction(
		ctx, id, requiredInMempool, timeout,
	); err != nil {
		return fmt.Errorf("error query tx failed: %w, broadcast initial error: %s", err, initErrMsg)
	}
	return nil
}
