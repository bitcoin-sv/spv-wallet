package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
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
// @Router		/v1/contacts [get]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	params := c.Request.URL.Query()

	queryParams, metadata, _, err := actions.GetSearchQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	dbConditions := make(map[string]interface{})

	for key, value := range params {
		dbConditions[key] = value
	}

	dbConditions["xpub_id"] = reqXPubID

	var contacts []*engine.Contact
	if contacts, err = a.Services.SpvWalletEngine.GetContacts(
		c.Request.Context(),
		metadata,
		&dbConditions,
		queryParams,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	contactContracts := make([]*models.Contact, 0)
	for _, contact := range contacts {
		contactContracts = append(contactContracts, mappings.MapToContactContract(contact))
	}

	c.JSON(http.StatusOK, contactContracts)
}
