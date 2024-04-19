package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
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
// @Param		page query int false "page"
// @Param		page_size query int false "page_size"
// @Param		order_by_field query string false "order_by_field"
// @Param		sort_direction query string false "sort_direction"
// @Param		conditions query string false "conditions"
// @Success		200 {object} []models.Contact "List of contacts"
// @Failure		400	"Bad request - Error while parsing SearchRequestParameters from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @Router		/v1/contact/search [POST]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	queryParams, metadata, conditions, err := actions.GetSearchQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	contacts, err := a.Services.SpvWalletEngine.GetContactsByXpubID(
		c.Request.Context(),
		reqXPubID,
		metadata,
		*conditions,
		queryParams,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, mappings.MapToContactContracts(contacts))
}
