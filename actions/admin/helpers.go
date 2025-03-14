package admin

import (
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// Helper function to prepare transaction query parameters
func prepareQueryParams(c *gin.Context, searchParams *filter.SearchParams[filter.AdminTransactionFilter]) *transactionQueryParams {
	return &transactionQueryParams{
		Context:     c.Request.Context(),
		XPubID:      searchParams.Conditions.XPubID,
		Metadata:    mappings.MapToMetadata(searchParams.Metadata),
		Conditions:  searchParams.Conditions.ToDbConditions(),
		PageOptions: mappings.MapToDbQueryParams(&searchParams.Page),
	}
}

// Helper function to fetch transactions based on XPubID presence
func fetchTransactions(c *gin.Context, params *transactionQueryParams) ([]*engine.Transaction, error) {
	if params.XPubID != nil {
		transactions, err := reqctx.Engine(c).GetTransactionsByXpubID(params.Context, *params.XPubID, params.Metadata, params.Conditions, params.PageOptions)
		if err != nil {
			return nil, fmt.Errorf("fetch transactions by XPubID failed: %w", err)
		}
		return transactions, nil
	}
	transactions, err := reqctx.Engine(c).GetTransactions(params.Context, params.Metadata, params.Conditions, params.PageOptions)
	if err != nil {
		return nil, fmt.Errorf("fetch transactions failed: %w", err)
	}
	return transactions, nil
}

// Helper function to count transactions based on XPubID presence
func countTransactions(c *gin.Context, params *transactionQueryParams) (int64, error) {
	if params.XPubID != nil {
		count, err := reqctx.Engine(c).GetTransactionsByXpubIDCount(params.Context, *params.XPubID, params.Metadata, params.Conditions)
		if err != nil {
			return 0, fmt.Errorf("count transactions by XPubID failed: %w", err)
		}
		return count, nil
	}

	count, err := reqctx.Engine(c).GetTransactionsCount(params.Context, params.Metadata, params.Conditions)
	if err != nil {
		return 0, fmt.Errorf("count transactions failed: %w", err)
	}
	return count, nil
}
