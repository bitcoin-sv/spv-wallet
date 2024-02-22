package chainstate

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func (c *Client) broadcastClientInit(ctx context.Context) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		defer txn.StartSegment("start_broadcast_client").End()
	}

	bc := c.options.config.broadcastClient
	if bc == nil {
		err := errors.New("broadcast client is not configured")
		return err
	}

	if c.isFeeQuotesEnabled() {
		// get the lowest fee
		var feeQuotes []*broadcast.FeeQuote
		feeQuotes, err := bc.GetFeeQuote(ctx)
		if err != nil {
			return err
		}
		if len(feeQuotes) == 0 {
			return errors.New("no fee quotes returned from broadcast client")
		}
		c.options.logger.Info().Msgf("got %d fee quote(s) from broadcast client", len(feeQuotes))
		fees := make([]utils.FeeUnit, len(feeQuotes))
		for index, fee := range feeQuotes {
			fees[index] = utils.FeeUnit{
				Satoshis: int(fee.MiningFee.Satoshis),
				Bytes:    int(fee.MiningFee.Bytes),
			}
		}
		c.options.config.feeUnit = utils.LowestFee(fees, c.options.config.feeUnit)
	}

	return nil
}
