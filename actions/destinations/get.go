package destinations

import (
	spverrors2 "github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// get will get an existing model
// Get Destination godoc
// @Summary		Get a destination
// @Description	Get a destination
// @Tags		Destinations
// @Produce		json
// @Param		id query string false "Destination ID"
// @Param		address query string false "Destination address"
// @Param		locking_script query string false "Destination locking script"
// @Success		200 {object} models.Destination "Destination with given id"
// @Failure		400	"Bad request - All parameters are missing (id, address, locking_script)"
// @Failure 	500	"Internal server error - Error while getting destination"
// @Router		/v1/destination [get]
// @Security	x-auth-xpub
func (a *Action) get(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	id := c.Query("id")
	address := c.Query("address")
	lockingScript := c.Query("locking_script")
	if id == "" && address == "" && lockingScript == "" {
		spverrors2.ErrorResponse(c, spverrors2.ErrOneOfTheFieldsIsRequired, a.Services.Logger)
		return
	}

	var destination *engine.Destination
	var err error
	if id != "" {
		destination, err = a.Services.SpvWalletEngine.GetDestinationByID(
			c.Request.Context(), reqXPubID, id,
		)
	} else if address != "" {
		destination, err = a.Services.SpvWalletEngine.GetDestinationByAddress(
			c.Request.Context(), reqXPubID, address,
		)
	} else {
		destination, err = a.Services.SpvWalletEngine.GetDestinationByLockingScript(
			c.Request.Context(), reqXPubID, lockingScript,
		)
	}
	if err != nil {
		spverrors2.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	contract := mappings.MapToDestinationContract(destination)
	c.JSON(http.StatusOK, contract)
}
