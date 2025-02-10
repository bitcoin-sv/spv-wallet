package v2

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/base"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/config"
)

// Server is the implementation of the server oapi-codegen's interface
type Server struct {
	admin.APIAdmin
	base.APIBase
}

// check if the Server implements the interface api.ServerInterface
var _ api.ServerInterface = &Server{}

// NewServer creates a new server
func NewServer(config *config.AppConfig) *Server {
	return &Server{
		APIAdmin: *admin.NewAPIAdmin(),
		APIBase:  *base.NewAPIBase(config),
	}
}
