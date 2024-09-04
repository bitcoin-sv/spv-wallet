package destinations

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// get will get an existing model
// Get Destination godoc
// @Summary		Get a destination. This endpoint has been deprecated (it will be removed in the future).
// @Description	Get a destination. This endpoint has been deprecated (it will be removed in the future).
// @Tags		Destinations
// @Produce		json
// @Param		id query string false "Destination ID"
// @Param		address query string false "Destination address"
// @Param		locking_script query string false "Destination locking script"
// @Success		200 {object} models.Destination "Destination with given id"
// @Failure		400	"Bad request - All parameters are missing (id, address, locking_script)"
// @Failure 	500	"Internal server error - Error while getting destination"
// @DeprecatedRouter  /v1/destination [get]
// @Security	x-auth-xpub
func get(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	engineInstance := reqctx.Engine(c)

	id := c.Query("id")
	address := c.Query("address")
	lockingScript := c.Query("locking_script")
	if id == "" && address == "" && lockingScript == "" {
		spverrors.ErrorResponse(c, spverrors.ErrOneOfTheFieldsIsRequired, logger)
		return
	}

	var destination *engine.Destination
	var err error
	if id != "" {
		destination, err = engineInstance.GetDestinationByID(
			c.Request.Context(), userContext.GetXPubID(), id,
		)
	} else if address != "" {
		destination, err = engineInstance.GetDestinationByAddress(
			c.Request.Context(), userContext.GetXPubID(), address,
		)
	} else {
		destination, err = engineInstance.GetDestinationByLockingScript(
			c.Request.Context(), userContext.GetXPubID(), lockingScript,
		)
	}
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contract := mappings.MapOldToDestinationContract(destination)
	c.JSON(http.StatusOK, contract)
}
