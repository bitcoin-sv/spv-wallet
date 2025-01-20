package actions

import (
	"github.com/bitcoin-sv/spv-wallet/api"
)

// Server is the implementation of the server oapi-codegen's interface
type Server struct {
	AdminServer
	UserServer
}

// check if the Server implements the interface api.ServerInterface
var _ api.ServerInterface = &Server{}

// NewServer creates a new server
func NewServer() Server {
	return Server{}
}
