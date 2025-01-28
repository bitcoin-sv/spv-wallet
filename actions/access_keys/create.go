package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// create will make a new model using the services defined in the action object
// Create access key godoc
// @Summary		Create access key
// @Description	Create access key
// @Tags		Access-key
// @Produce		json
// @Param		CreateAccessKey body CreateAccessKey true " "
// @Success		201	{object} response.AccessKey "Created AccessKey"
// @Failure		400	"Bad request - Error while parsing CreateAccessKey from request body"
// @Failure 	500	"Internal server error - Error while creating new access key"
// @Router		/api/v1/users/current/keys [post]
// @Security	x-auth-xpub
func create(c *gin.Context, userContext *reqctx.UserContext) {
	xpub, err := userContext.ShouldGetXPub()
	if err != nil {
		spverrors.AbortWithErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	logger := reqctx.Logger(c)
	var requestBody CreateAccessKey
	if err := c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	// Create a new accessKey
	accessKey, err := reqctx.Engine(c).NewAccessKey(
		c.Request.Context(),
		xpub,
		engine.WithMetadatas(requestBody.Metadata),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contract := mappings.MapToAccessKeyContract(accessKey)
	c.JSON(http.StatusCreated, contract)
}
