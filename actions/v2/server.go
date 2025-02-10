package v2

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/data"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/operations"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/transactions"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/users"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/rs/zerolog"
)

// Server is the implementation of the server oapi-codegen's interface
type Server struct {
	admin.Server
	data.APIData
	users.APIUsers
	operations.APIOperations
	transactions.APITransactions
}

// check if the Server implements the interface api.ServerInterface
var _ api.ServerInterface = &Server{}

// NewServer creates a new server
func NewServer(config *config.AppConfig, engine engine.ClientInterface, logger *zerolog.Logger) *Server {
	return &Server{
		admin.Server{},
		data.NewAPIData(engine, logger),
		users.NewAPIUsers(engine, logger),
		operations.NewAPIOperations(engine, logger),
		transactions.NewAPITransactions(engine, logger),
	}
}
