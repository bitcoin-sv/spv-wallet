package destinations

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// update will update an existing model
// Update Destination godoc
// @Summary		Update destination. This endpoint has been deprecated (it will be removed in the future).
// @Description	Update destination. This endpoint has been deprecated (it will be removed in the future).
// @Tags		Destinations
// @Produce		json
// @Param		UpdateDestination body UpdateDestination false " "
// @Success		200 {object} models.Destination "Updated Destination"
// @Failure		400	"Bad request - Error while parsing UpdateDestination from request body"
// @Failure 	500	"Internal Server Error - Error while updating destination"
// @DeprecatedRouter  /v1/destination [patch]
// @Security	x-auth-xpub
func update(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	engineInstance := reqctx.Engine(c)

	var requestBody UpdateDestination
	if err := c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}
	if requestBody.ID == "" && requestBody.Address == "" && requestBody.LockingScript == "" {
		spverrors.ErrorResponse(c, spverrors.ErrOneOfTheFieldsIsRequired, logger)
		return
	}

	// Get the destination
	var destination *engine.Destination
	var err error
	if requestBody.ID != "" {
		destination, err = engineInstance.UpdateDestinationMetadataByID(
			c.Request.Context(), userContext.GetXPubID(), requestBody.ID, requestBody.Metadata,
		)
	} else if requestBody.Address != "" {
		destination, err = engineInstance.UpdateDestinationMetadataByAddress(
			c.Request.Context(), userContext.GetXPubID(), requestBody.Address, requestBody.Metadata,
		)
	} else {
		destination, err = engineInstance.UpdateDestinationMetadataByLockingScript(
			c.Request.Context(), userContext.GetXPubID(), requestBody.LockingScript, requestBody.Metadata,
		)
	}
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contract := mappings.MapOldToDestinationContract(destination)
	c.JSON(http.StatusOK, contract)
}
