package users

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// oldGet will get an existing model
// Get current user information godoc
// @Summary		Get current user information - Use (GET) /api/v1/users/current instead.
// @Description	This endpoint has been deprecated. Use (GET) /api/v1/users/current instead.
// @Tags		Users
// @Produce		json
// @Success		200 {object} models.Xpub "xPub associated with the given xPub from auth header"
// @Failure		500	"Internal Server Error - Error while fetching xPub"
// @DeprecatedRouter  /v1/xpub [get]
// @Security	x-auth-xpub
func oldGet(c *gin.Context, userContext *reqctx.UserContext) {
	getHelper(c, true, userContext)
}

// get will get an existing model
// Get current user information godoc
// @Summary		Get current user information
// @Description	Get current user information
// @Tags		Users
// @Produce		json
// @Success		200 {object} response.Xpub "xPub associated with the given xPub from auth header"
// @Failure		500	"Internal Server Error - Error while fetching xPub"
// @Router		/api/v1/users/current [get]
// @Security	x-auth-xpub
func get(c *gin.Context, userContext *reqctx.UserContext) {
	getHelper(c, false, userContext)
}

func getHelper(c *gin.Context, snakeCase bool, userContext *reqctx.UserContext) {
	if snakeCase {
		contract := mappings.MapToOldXpubContract(userContext.GetXPubObj())
		c.JSON(http.StatusOK, contract)
		return
	}

	contract := mappings.MapToXpubContract(userContext.GetXPubObj())
	c.JSON(http.StatusOK, contract)
}
