package engine

import (
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/record"
)

// TransactionOutlinesService will return the outlines.Service if it exists
func (c *Client) TransactionOutlinesService() outlines.Service {
	return c.options.transactionOutlinesService
}

// TransactionRecordService will return the record.Service if it exists
func (c *Client) TransactionRecordService() *record.Service {
	return c.options.transactionRecordService
}
