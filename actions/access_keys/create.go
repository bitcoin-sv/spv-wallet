package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// create will make a new model using the services defined in the action object
// Create access key godoc
// @Summary		Create access key
// @Description	Create access key
// @Tags		Access-key
// @Produce		json
// @Param		CreateAccessKey body CreateAccessKey true "CreateAccessKey model containing metadata"
// @Success		201
// @Router		/v1/access-key [post]
// @Security	x-auth-xpub
func (a *Action) create(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)

	var requestBody CreateAccessKey
	if err := c.Bind(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Create a new accessKey
	accessKey, err := a.Services.SpvWalletEngine.NewAccessKey(
		c.Request.Context(),
		reqXPub,
		engine.WithMetadatas(requestBody.Metadata),
	)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	contract := mappings.MapToAccessKeyContract(accessKey)
	c.JSON(http.StatusCreated, contract)
}
