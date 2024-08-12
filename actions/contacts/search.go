package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// Search will fetch a list of contacts
// Get contacts godoc
// @Summary		Search contacts - Use (GET) /api/v1/contacts instead.
// @Description	This endpoint has been deprecated. Use (GET) /api/v1/contacts instead.
// @Tags		Contact
// @Produce		json
// @Param		SearchContacts body filter.SearchContacts false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} models.SearchContactsResponse "List of contacts"
// @Failure		400	"Bad request - Error while parsing SearchContacts from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @DeprecatedRouter  /v1/contact/search [post]
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

	contracts := mappings.MapToOldContactContracts(contacts)

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
// @Success		200 {object} response.PageModel[response.Contact] "Page of contacts"
// @Failure		400	"Bad request - Error while parsing SearchContacts from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @Router		/api/v1/contacts [get]
// @Security	x-auth-xpub
func (a *Action) getContacts(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	contacts, count := a.searchContacts(c, reqXPubID, "")
	if contacts == nil {
		return
	}

	contracts := mappings.MapToContactContracts(contacts)

	totalPages := 0
	if int(count) != 0 {
		totalPages = len(contracts) / int(count)
	}

	response := response.PageModel[response.Contact]{
		Content: contracts,
		Page: response.PageDescription{
			Size:          len(contracts),
			Number:        0,
			TotalElements: int(count),
			TotalPages:    totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}

// getContactByPaymail will fetch a list of contacts
// @Summary		Get contact by paymail
// @Description	Get contact by paymail
// @Tags		Contacts
// @Produce		json
// @Success		200 {object} response.Contact "Contact"
// @Failure		400	"Bad request - Error while parsing SearchContacts from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @Router		/api/v1/contacts/{paymail} [get]
// @Security	x-auth-xpub
func (a *Action) getContactByPaymail(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	paymail := c.Param("paymail")

	contacts, _ := a.searchContacts(c, reqXPubID, paymail)
	if contacts == nil {
		return
	}

	if contacts == nil || len(contacts) != 1 {
		spverrors.ErrorResponse(c, spverrors.ErrContactNotFound, a.Services.Logger)
		return
	}

	contracts := mappings.MapToContactContracts(contacts)

	c.JSON(http.StatusOK, contracts[0])
}

// searchContacts - a helper function for searching contacts
func (a *Action) searchContacts(c *gin.Context, reqXPubID string, paymail string) ([]*engine.Contact, int64) {
	var reqParams filter.SearchContacts

	if paymail != "" {
		preConditions := filter.ContactFilter{Paymail: &paymail}
		reqParams.ConditionsModel = filter.ConditionsModel[filter.ContactFilter]{
			Conditions: &preConditions,
		}
	}

	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return nil, 0
	}

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidConditions, a.Services.Logger)
		return nil, 0
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
		return nil, 0
	}

	count, err := a.Services.SpvWalletEngine.GetContactsByXPubIDCount(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return nil, 0
	}

	return contacts, count
}
