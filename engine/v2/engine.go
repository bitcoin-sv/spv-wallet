package v2

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/addresses"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/data"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database/repository"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/operations"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/record"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users"
)

// Engine represents the engine of the SPV wallet
type Engine interface {
	Repositories() *repository.All
	UsersService() *users.Service
	PaymailsService() *paymails.Service
	AddressesService() *addresses.Service
	DataService() *data.Service
	OperationsService() *operations.Service
	TransactionOutlinesService() outlines.Service
	TransactionRecordService() *record.Service
	Close(ctx context.Context) error
}
