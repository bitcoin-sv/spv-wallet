package chainstate

import (
	"context"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

func (c *Client) broadcastClientInit(ctx context.Context) error {

	bc := c.options.config.broadcastClient
	if bc == nil {
		err := spverrors.Newf("broadcast client is not configured")
		return err
	}

	if c.isFeeQuotesEnabled() {
		// get the lowest fee
		var feeQuotes []*broadcast.FeeQuote
		feeQuotes, err := bc.GetFeeQuote(ctx)
		if err != nil {
			return spverrors.Wrapf(err, "failed to get fee quotes from broadcast client")
		}
		if len(feeQuotes) == 0 {
			return spverrors.Newf("no fee quotes returned from broadcast client")
		}
		c.options.logger.Info().Msgf("got %d fee quote(s) from broadcast client", len(feeQuotes))
		fees := make([]bsv.FeeUnit, len(feeQuotes))
		for index, fee := range feeQuotes {
			fees[index] = bsv.FeeUnit{
				Satoshis: int(fee.MiningFee.Satoshis),
				Bytes:    int(fee.MiningFee.Bytes),
			}
		}
		c.options.config.feeUnit = utils.LowestFee(fees, c.options.config.feeUnit)
	}

	return nil
}
