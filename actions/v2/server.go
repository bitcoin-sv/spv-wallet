package v2

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/rs/zerolog"
)

// Server is the implementation of the server oapi-codegen's interface
type Server struct {
	admin.APIAdmin
}

// check if the Server implements the interface api.ServerInterface
var _ api.ServerInterface = &Server{}

// NewServer creates a new server
func NewServer(config *config.AppConfig, engine engine.ClientInterface, logger *zerolog.Logger) *Server {
	return &Server{
		APIAdmin: *admin.NewAPIAdmin(engine, logger),
	}
}
