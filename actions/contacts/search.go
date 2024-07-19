package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// Search will fetch a list of contacts
// Get contacts godoc
// @Summary		Search contacts
// @Description	Search contacts
// @Tags		Contact
// @Produce		json
// @Param		SearchContacts body filter.SearchContacts false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} models.SearchContactsResponse "List of contacts"
// @Failure		400	"Bad request - Error while parsing SearchContacts from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @Router		/v1/contact/search [POST]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams filter.SearchContacts
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidConditions, a.Services.Logger)
		return
	}

	reqParams.DefaultsIfNil()

	contacts, err := a.Services.SpvWalletEngine.GetContactsByXpubID(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	contracts := mappings.MapToContactContracts(contacts)

	count, err := a.Services.SpvWalletEngine.GetContactsByXPubIDCount(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	response := models.SearchContactsResponse{
		Content: contracts,
		Page:    common.GetPageFromQueryParams(reqParams.QueryParams, count),
	}

	c.JSON(http.StatusOK, response)
}

// TODO: this method will be changed based on search poc
// getContacts will fetch a list of contacts
// @Summary		Get contacts
// @Description	Get contacts
// @Tags		Contacts
// @Produce		json
// @Param		SearchContacts body filter.SearchContacts false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} models.PageModel[models.Contact] "Page of contacts"
// @Failure		400	"Bad request - Error while parsing SearchContacts from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @Router		/v1/contacts/ [GET]
// @Security	x-auth-xpub
func (a *Action) getContacts(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams filter.SearchContacts
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	reqParams.DefaultsIfNil()

	contacts, err := a.Services.SpvWalletEngine.GetContactsByXpubID(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contracts := mappings.MapToContactContracts(contacts)

	count, err := a.Services.SpvWalletEngine.GetContactsByXPubIDCount(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := models.PageModel[models.Contact]{
		Content: contracts,
		Page: models.PageDescription{
			Size:          len(contracts),
			Number:        0,
			TotalElements: int(count),
			TotalPages:    len(contracts) / int(count),
		},
	}

	c.JSON(http.StatusOK, response)
}

// getContactByPaymail will fetch a list of contacts
// @Summary		Get contact by paymail
// @Description	Get contact by paymail
// @Tags		Contacts
// @Produce		json
// @Success		200 {object} models.Contact "Contact"
// @Failure		400	"Bad request - Error while parsing SearchContacts from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @Router		/v1/contacts/{paymail} [GET]
// @Security	x-auth-xpub
func (a *Action) getContactByPaymail(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	paymail := c.Param("paymail")
	var reqParams filter.SearchContacts
	reqParams.ConditionsModel.Conditions.Paymail = &paymail

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	reqParams.DefaultsIfNil()

	contacts, err := a.Services.SpvWalletEngine.GetContactsByXpubID(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if contacts == nil || len(contacts) != 1 {
		c.JSON(http.StatusNotFound, "contact not found")
		return
	}

	c.JSON(http.StatusOK, contacts[0])
}
