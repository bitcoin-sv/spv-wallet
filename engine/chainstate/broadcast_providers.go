package chainstate

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

// generic broadcast provider
type txBroadcastProvider interface {
	getName() string
	broadcast(ctx context.Context, c *Client) *BroadcastFailure
}

// BroadcastClient provider
type broadcastClientProvider struct {
	txID, txHex string
	format      HexFormatFlag
}

func (provider *broadcastClientProvider) getName() string {
	return ProviderBroadcastClient
}

func (provider *broadcastClientProvider) broadcast(ctx context.Context, c *Client) *BroadcastFailure {
	logger := c.options.logger

	logger.Debug().
		Str("txID", provider.txID).
		Msgf("executing broadcast request for %s", provider.getName())

	tx := broadcast.Transaction{
		Hex: provider.txHex,
	}

	formatOpt := broadcast.WithRawFormat()
	if provider.format.Contains(Ef) {
		formatOpt = broadcast.WithEfFormat()

		logger.Debug().
			Str("txID", provider.txID).
			Msgf("broadcast with broadcast-client in Extended Format")
	} else {
		logger.Debug().
			Str("txID", provider.txID).
			Msgf("broadcast with broadcast-client in RawTx format")
	}

	result, err := c.BroadcastClient().SubmitTransaction(
		ctx,
		&tx,
		formatOpt,
		broadcast.WithCallback(c.options.config.callbackURL, c.options.config.callbackToken),
	)

	if err != nil {
		var arcError *broadcast.ArcError
		if errors.As(err, &arcError) {
			logger.Debug().
				Str("txID", provider.txID).
				Msgf("error broadcast request for %s failed: %s", provider.getName(), arcError.Error())

			return &BroadcastFailure{
				InvalidTx: arcError.IsRejectedTransaction(),
				Error:     arcError,
			}
		}

		return &BroadcastFailure{
			InvalidTx: false,
			Error:     err,
		}
	}

	logger.Debug().
		Str("txID", provider.txID).
		Msgf("result broadcast request for %s blockhash: %s status: %s", provider.getName(), result.BlockHash, result.TxStatus.String())

	return nil
}
