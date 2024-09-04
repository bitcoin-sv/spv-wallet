package destinations

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// create will make a new destination
// Count Destinations godoc
// @Summary		Create a new destination. This endpoint has been deprecated (it will be removed in the future).
// @Description	Create a new destination. This endpoint has been deprecated (it will be removed in the future).
// @Tags		Destinations
// @Produce		json
// @Param		CreateDestination body CreateDestination false " "
// @Success		201 {object} models.Destination "Created Destination"
// @Failure		400	"Bad request - Error while parsing CreateDestination from request body"
// @Failure 	500	"Internal Server Error - Error while creating destination"
// @DeprecatedRouter  /v1/destination [post]
// @Security	x-auth-xpub
func create(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	engineInstance := reqctx.Engine(c)

	var requestBody CreateDestination
	err := c.Bind(&requestBody)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	opts := engineInstance.DefaultModelOptions()

	if requestBody.Metadata != nil {
		opts = append(opts, engine.WithMetadatas(requestBody.Metadata))
	}

	var destination *engine.Destination
	if destination, err = engineInstance.NewDestination(
		c.Request.Context(),
		userContext.GetXPub(),
		uint32(0), // todo: use a constant? protect this?
		utils.ScriptTypePubKeyHash,
		opts...,
	); err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contract := mappings.MapOldToDestinationContract(destination)
	c.JSON(http.StatusCreated, contract)
}
