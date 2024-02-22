package chainstate

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/libsv/go-bc"
	"github.com/tonicpow/go-minercraft/v2"
)

// query will try ALL providers in order and return the first "valid" response based on requirements
func (c *Client) query(ctx context.Context, id string, requiredIn RequiredIn,
	timeout time.Duration,
) *TransactionInfo {
	// Create a context (to cancel or timeout)
	ctxWithCancel, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	switch c.ActiveProvider() {
	case ProviderMinercraft:
		for index := range c.options.config.minercraftConfig.queryMiners {
			if c.options.config.minercraftConfig.queryMiners[index] != nil {
				if res, err := queryMinercraft(
					ctxWithCancel, c, c.options.config.minercraftConfig.queryMiners[index], id,
				); err == nil && checkRequirementMapi(requiredIn, id, res) {
					return res
				}
			}
		}
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
	case ProviderMinercraft:
		for index := range c.options.config.minercraftConfig.queryMiners {
			wg.Add(1)
			go func(
				ctx context.Context, client *Client,
				wg *sync.WaitGroup, miner *minercraft.Miner,
				id string, requiredIn RequiredIn,
			) {
				defer wg.Done()
				if res, err := queryMinercraft(
					ctx, client, miner, id,
				); err == nil && checkRequirementMapi(requiredIn, id, res) {
					resultsChannel <- res
				}
			}(ctxWithCancel, c, &wg, c.options.config.minercraftConfig.queryMiners[index], id, requiredIn)
		}
	case ProviderBroadcastClient:
		wg.Add(1)
		go func(ctx context.Context, client *Client, id string, requiredIn RequiredIn) {
			defer wg.Done()
			if resp, err := queryBroadcastClient(
				ctx, client, id,
			); err == nil && checkRequirementArc(requiredIn, id, resp) {
				resultsChannel <- resp
			}
		}(ctxWithCancel, c, id, requiredIn)
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

// queryMinercraft will submit a query transaction request to a miner using Minercraft(mAPI or Arc)
func queryMinercraft(ctx context.Context, client ClientInterface, miner *minercraft.Miner, id string) (*TransactionInfo, error) {
	client.DebugLog("executing request in minercraft using miner: " + miner.Name)
	if resp, err := client.Minercraft().QueryTransaction(ctx, miner, id, minercraft.WithQueryMerkleProof()); err != nil {
		client.DebugLog("error executing request in minercraft using miner: " + miner.Name + " failed: " + err.Error())
		return nil, err
	} else if resp != nil && resp.Query.ReturnResult == mAPISuccess && strings.EqualFold(resp.Query.TxID, id) {
		return &TransactionInfo{
			BlockHash:     resp.Query.BlockHash,
			BlockHeight:   resp.Query.BlockHeight,
			Confirmations: resp.Query.Confirmations,
			ID:            resp.Query.TxID,
			MinerID:       resp.Query.MinerID,
			Provider:      miner.Name,
			MerkleProof:   resp.Query.MerkleProof,
		}, nil
	}
	return nil, ErrTransactionIDMismatch
}

// queryBroadcastClient will submit a query transaction request to a go-broadcast-client
func queryBroadcastClient(ctx context.Context, client ClientInterface, id string) (*TransactionInfo, error) {
	client.DebugLog("executing request using " + ProviderBroadcastClient)
	if resp, err := client.BroadcastClient().QueryTransaction(ctx, id); err != nil {
		client.DebugLog("error executing request using " + ProviderBroadcastClient + " failed: " + err.Error())
		return nil, err
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
