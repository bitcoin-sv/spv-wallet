package chainstate

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/tonicpow/go-minercraft/v2"
)

// generic broadcast provider
type txBroadcastProvider interface {
	getName() string
	broadcast(ctx context.Context, c *Client) error
}

// mAPI provider
type mapiBroadcastProvider struct {
	miner       *minercraft.Miner
	txID, txHex string
}

func (provider mapiBroadcastProvider) getName() string {
	return provider.miner.Name
}

func (provider mapiBroadcastProvider) broadcast(ctx context.Context, c *Client) error {
	return broadcastMAPI(ctx, c, provider.miner, provider.txID, provider.txHex)
}

// broadcastMAPI will broadcast a transaction to a miner using mAPI
func broadcastMAPI(ctx context.Context, client ClientInterface, miner *minercraft.Miner, id, hex string) error {
	debugLog(client, id, "executing broadcast request in mapi using miner: "+miner.Name)

	resp, err := client.Minercraft().SubmitTransaction(ctx, miner, &minercraft.Transaction{
		CallBackEncryption: "", // todo: allow customizing the payload
		CallBackToken:      "",
		CallBackURL:        "",
		DsCheck:            false,
		MerkleFormat:       "",
		MerkleProof:        false,
		RawTx:              hex,
	})
	if err != nil {
		debugLog(client, id, "error executing request in mapi using miner: "+miner.Name+" failed: "+err.Error())
		return err
	}

	// Something went wrong - got back an id that does not match
	if resp == nil {
		return emptyBroadcastResponseErr(id)
	}
	if !strings.EqualFold(resp.Results.TxID, id) {
		return incorrectTxIDReturnedErr(resp.Results.TxID, id)
	}

	// mAPI success of broadcast
	if resp.Results.ReturnResult == mAPISuccess {
		return nil
	}

	// Check error message (for success error message)
	if doesErrorContain(resp.Results.ResultDescription, broadcastSuccessErrors) {
		return nil
	}

	// We got a potential real error message?
	return errors.New(resp.Results.ResultDescription)
}

func incorrectTxIDReturnedErr(actualTxID, expectedTxID string) error {
	return fmt.Errorf("returned tx id [%s] does not match given tx id [%s]", actualTxID, expectedTxID)
}

func emptyBroadcastResponseErr(txID string) error {
	return fmt.Errorf("an empty response was returned after broadcasting of tx id [%s]", txID)
}

////

// BroadcastClient provider
type broadcastClientProvider struct {
	txID, txHex string
}

func (provider broadcastClientProvider) getName() string {
	return ProviderBroadcastClient
}

// Broadcast using BroadcastClient
func (provider broadcastClientProvider) broadcast(ctx context.Context, c *Client) error {
	return broadcastWithBroadcastClient(ctx, c, provider.txID, provider.txHex)
}

func broadcastWithBroadcastClient(ctx context.Context, client *Client, txID, hex string) error {
	debugLog(client, txID, "executing broadcast request for "+ProviderBroadcastClient)

	tx := broadcast.Transaction{
		Hex: hex,
	}

	result, err := client.BroadcastClient().SubmitTransaction(
		ctx,
		&tx,
		broadcast.WithRawFormat(),
		broadcast.WithCallback(client.options.config.callbackURL, client.options.config.callbackToken),
	)
	if err != nil {
		debugLog(client, txID, "error broadcast request for "+ProviderBroadcastClient+" failed: "+err.Error())
		return err
	}

	debugLog(client, txID, "result broadcast request for "+ProviderBroadcastClient+" blockhash: "+result.BlockHash+" status: "+result.TxStatus.String())

	return nil
}
