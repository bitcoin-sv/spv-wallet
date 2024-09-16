package engine

import (
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
)

// TransactionDraftService will return the draft.Service if it exists
func (c *Client) TransactionDraftService() draft.Service {
	if c.options.transactionDraftService != nil {
		return c.options.transactionDraftService
	}
	return nil
}
