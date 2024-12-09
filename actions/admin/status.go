package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// @Summary				Get status Use (GET) /api/v1/admin/status instead.
// @Description			This endpoint has been deprecated. Use (GET) /api/v1/admin/status instead.
// @Tags				Admin
// @Produce				json
// @Success				200 {boolean} bool "Status response"
// @DeprecatedRouter 	/v1/admin/status [get]
// @Security			x-auth-xpub
func statusOld(c *gin.Context, _ *reqctx.AdminContext) {
	c.JSON(http.StatusOK, true)
}

// @Summary			Get status
// @Description		Get status
// @Tags			Admin
// @Produce			json
// @Success			200 {boolean} bool "Status response"
// @Router			/api/v1/admin/status [get]
// @Security		x-auth-xpub
func status(c *gin.Context, _ *reqctx.AdminContext) {
	c.JSON(http.StatusOK, true)
}
