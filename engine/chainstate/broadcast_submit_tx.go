package chainstate

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func (chainstate *Client) submitTransaction(ctx context.Context, txID, txHex string, format HexFormatFlag) *BroadcastFailure {
	logger := chainstate.options.logger

	logger.Debug().
		Str("txID", txID).
		Msgf("executing broadcast request")

	tx := broadcast.Transaction{
		Hex: txHex,
	}

	formatOpt := broadcast.WithRawFormat()
	if format.Contains(Ef) {
		formatOpt = broadcast.WithEfFormat()

		logger.Debug().
			Msgf("broadcast with broadcast-client in Extended Format")
	} else {
		logger.Debug().
			Str("txID", txID).
			Msgf("broadcast with broadcast-client in RawTx format")
	}

	result, err := chainstate.BroadcastClient().SubmitTransaction(
		ctx,
		&tx,
		formatOpt,
		broadcast.WithCallback(chainstate.options.config.callbackURL, chainstate.options.config.callbackToken),
	)
	if err != nil {
		var arcError *broadcast.ArcError
		if errors.As(err, &arcError) {
			logger.Debug().
				Str("txID", txID).
				Msgf("broadcast failed: %s", arcError.Error())

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
		Str("txID", txID).
		Msgf("result broadcast request - blockhash: %s status: %s", result.BlockHash, result.TxStatus.String())

	return nil
}
