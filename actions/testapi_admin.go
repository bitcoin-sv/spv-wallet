package actions

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/gin-gonic/gin"
)

// AdminServer is the partial implementation of the server oapi-codegen's interface
// NOTE: This is showcase how to split Server implementation into multiple packages/structs/files
type AdminServer struct{}

// GetAdminTestapi is the implementation of the corresponding endpoint
func (AdminServer) GetAdminTestapi(c *gin.Context) {
	c.JSON(200, api.ExResponse{
		XpubID:    "admin",
		Something: "something",
	})
}
