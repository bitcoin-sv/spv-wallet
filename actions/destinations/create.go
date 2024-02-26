package destinations

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// create will make a new destination
// Count Destinations godoc
// @Summary		Create a new destination
// @Description	Create a new destination
// @Tags		Destinations
// @Produce		json
// @Param		type query string false "type"
// @Param		reference_id query string false "reference_id"
// @Param		metadata query string false "metadata"
// @Success		201
// @Router		/v1/destination [post]
// @Security	x-auth-xpub
func (a *Action) create(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)
	xPub, err := a.Services.SpvWalletEngine.GetXpub(c.Request.Context(), reqXPub)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	} else if xPub == nil {
		c.JSON(http.StatusForbidden, actions.ErrXpubNotFound)
		return
	}

	var requestBody CreateDestination
	if err = c.Bind(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	opts := a.Services.SpvWalletEngine.DefaultModelOptions()

	if requestBody.Metadata != nil {
		opts = append(opts, engine.WithMetadatas(requestBody.Metadata))
	}

	var destination *engine.Destination
	if destination, err = a.Services.SpvWalletEngine.NewDestination(
		c.Request.Context(),
		xPub.RawXpub(),
		uint32(0), // todo: use a constant? protect this?
		utils.ScriptTypePubKeyHash,
		opts...,
	); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	contract := mappings.MapToDestinationContract(destination)
	c.JSON(http.StatusCreated, contract)
}
