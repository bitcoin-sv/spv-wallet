package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// create will make a new model using the services defined in the action object
// Create access key godoc
// @Summary		Create access key - Use (POST) /api/v1/users/current/keys instead.
// @Description	This endpoint has been deprecated. Use (POST) /api/v1/users/current/keys instead.
// @Tags		Access-key
// @Produce		json
// @Param		CreateAccessKey body CreateAccessKey true " "
// @Success		201	{object} models.AccessKey "Created AccessKey"
// @Failure		400	"Bad request - Error while parsing CreateAccessKey from request body"
// @Failure 	500	"Internal server error - Error while creating new access key"
// @DeprecatedRouter  /v1/access-key [post]
// @Security	x-auth-xpub
func (a *Action) oldCreate(c *gin.Context) {
	a.createHelper(c, true)
}

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
func (a *Action) create(c *gin.Context) {
	a.createHelper(c, false)
}

func (a *Action) createHelper(c *gin.Context, snakeCase bool) {
	reqXPub := c.GetString(auth.ParamXPubKey)

	var requestBody CreateAccessKey
	if err := c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	// Create a new accessKey
	accessKey, err := a.Services.SpvWalletEngine.NewAccessKey(
		c.Request.Context(),
		reqXPub,
		engine.WithMetadatas(requestBody.Metadata),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	if snakeCase {
		contract := mappings.MapToOldAccessKeyContract(accessKey)
		c.JSON(http.StatusCreated, contract)
		return
	}

	contract := mappings.MapToAccessKeyContract(accessKey)
	c.JSON(http.StatusCreated, contract)
}
