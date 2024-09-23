package chainstate

import (
	"context"
	"strings"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/libsv/go-bc"
)

// query will try ALL providers in order and return the first "valid" response based on requirements
func (c *Client) query(ctx context.Context, id string, requiredIn RequiredIn,
	timeout time.Duration,
) *TransactionInfo {
	// Create a context (to cancel or timeout)
	ctxWithCancel, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp, err := queryBroadcastClient(
		ctxWithCancel, c, id,
	)
	if err == nil && checkRequirementArc(requiredIn, id, resp) {
		return resp
	}

	return nil // No transaction information found
}

// queryBroadcastClient will submit a query transaction request to a go-broadcast-client
func queryBroadcastClient(ctx context.Context, client ClientInterface, id string) (*TransactionInfo, error) {
	client.DebugLog("executing request using " + ProviderBroadcastClient)
	if resp, failure := client.BroadcastClient().QueryTransaction(ctx, id); failure != nil {
		client.DebugLog("error executing request using " + ProviderBroadcastClient + " failed: " + failure.Error())
		return nil, spverrors.Wrapf(failure, "failed to query transaction using %s", ProviderBroadcastClient)
	} else if resp != nil && strings.EqualFold(resp.TxID, id) {
		bump, err := bc.NewBUMPFromStr(resp.BaseTxResponse.MerklePath)
		if err != nil {
			return nil, spverrors.Wrapf(err, "failed to parse BUMP from response: %s", resp.BaseTxResponse.MerklePath)
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
	return nil, spverrors.ErrTransactionIDMismatch
}
