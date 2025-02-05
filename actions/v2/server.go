package v2

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin"
	"github.com/bitcoin-sv/spv-wallet/api"
)

// Server is the implementation of the server oapi-codegen's interface
type Server struct {
	admin.AdminServer
}

// check if the Server implements the interface api.ServerInterface
var _ api.ServerInterface = &Server{}

// NewServer creates a new server
func NewServer() *Server {
	return &Server{}
}
