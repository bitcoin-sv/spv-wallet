package v2

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/gin-gonic/gin"
)

// Server is the implementation of the server oapi-codegen's interface
type Server struct {
	admin.Server
}

// check if the Server implements the interface api.ServerInterface
var _ api.ServerInterface = &Server{}

// NewServer creates a new server
func NewServer() *Server {
	return &Server{}
}

// GetApiV2Status is the handler for the status endpoint which is available for everyone
func (s *Server) GetApiV2Status(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
