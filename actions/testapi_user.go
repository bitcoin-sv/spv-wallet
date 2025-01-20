package actions

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// UserServer is the implementation of the server oapi-codegen's interface
type UserServer struct{}

// GetTestapi is a handler for the GET /testapi endpoint.
func (UserServer) GetTestapi(c *gin.Context, request api.GetTestapiParams) {
	user := reqctx.GetUserContext(c)

	c.JSON(200, api.ExResponse{
		XpubID:    user.GetXPubID(),
		Something: request.Something,
	})
}
