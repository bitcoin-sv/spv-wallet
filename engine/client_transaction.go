package engine

import (
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
)

// TransactionOutlinesService will return the outlines.Service if it exists
func (c *Client) TransactionOutlinesService() outlines.Service {
	if c.options.transactionOutlinesService != nil {
		return c.options.transactionOutlinesService
	}
	return nil
}
