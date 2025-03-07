package transactions

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/rs/zerolog"
)

type transactionsRecordService interface {
	RecordTransactionOutline(ctx context.Context, userID string, outline *outlines.Transaction) (*txmodels.RecordedOutline, error)
}

// APIAdminTransactions represents server with admin API endpoints
type APIAdminTransactions struct {
	transactionsRecordService transactionsRecordService
	logger                    *zerolog.Logger
}

// NewAPIAdminTransactions creates a new APIAdminTransactions
func NewAPIAdminTransactions(engine engine.ClientInterface, logger *zerolog.Logger) APIAdminTransactions {
	return APIAdminTransactions{
		transactionsRecordService: engine.TransactionRecordService(),
		logger:                    logger,
	}
}
