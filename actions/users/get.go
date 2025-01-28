package users

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

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
	contract := mappings.MapToXpubContract(userContext.GetXPubObj())
	c.JSON(http.StatusOK, contract)
}
