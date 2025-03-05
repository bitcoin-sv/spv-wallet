package v2

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/base"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/data"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/merkleroots"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/operations"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/transactions"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/users"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/config"
	v2 "github.com/bitcoin-sv/spv-wallet/engine/v2"
	"github.com/rs/zerolog"
)

// apiV2 is the implementation of the server oapi-codegen's interface
type apiV2 struct {
	admin.APIAdmin
	base.APIBase
	data.APIData
	users.APIUsers
	operations.APIOperations
	transactions.APITransactions
	merkleroots.APIMerkleRoots
}

// NewV2API creates a new server
func NewV2API(config *config.AppConfig, engine v2.Engine, logger *zerolog.Logger) api.ServerInterface {
	return &apiV2{
		admin.NewAPIAdmin(engine, logger),
		base.NewAPIBase(config),
		data.NewAPIData(engine, logger),
		users.NewAPIUsers(engine, logger),
		operations.NewAPIOperations(engine, logger),
		transactions.NewAPITransactions(engine, logger),
		merkleroots.NewAPIMerkleRoots(engine, logger),
	}
}
