package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// Search will fetch a list of contacts
// Get contacts godoc
// @Summary		Search contacts
// @Description	Search contacts
// @Tags		Contact
// @Produce		json
// @Param		SearchContacts body SearchContacts false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Contact "List of contacts"
// @Failure		400	"Bad request - Error while parsing SearchContacts from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @Router		/v1/contact/search [POST]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams SearchContacts
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	contacts, err := a.Services.SpvWalletEngine.GetContactsByXpubID(
		c.Request.Context(),
		reqXPubID,
		reqParams.Metadata,
		conditions,
		reqParams.QueryParams,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, mappings.MapToContactContracts(contacts))
}
