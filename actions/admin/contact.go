package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// contactsSearch will fetch a list of contacts filtered by Metadata and AdminContactFilters
// Search for contacts filtering by metadata and AdminContactFilters godoc
// @Summary		Search for contacts
// @Description	Search for contacts
// @Tags		Admin
// @Produce		json
// @Param		AdminSearchContacts body filter.AdminContactFilter false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} response.PageModel[response.Contact] "List of contacts"
// @Failure		400	"Bad request - Error while parsing AdminSearchContacts from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @Router		/api/v1/admin/contacts [get]
// @Security	x-auth-xpub
func contactsSearch(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	engine := reqctx.Engine(c)

	searchParams, err := query.ParseSearchParams[filter.AdminContactFilter](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams.WithTrace(err), logger)
		return
	}

	conditions, err := searchParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidConditions.WithTrace(err), logger)
		return
	}
	metadata := mappings.MapToMetadata(searchParams.Metadata)
	pageOptions := mappings.MapToDbQueryParams(&searchParams.Page)

	contacts, err := engine.GetContacts(
		c.Request.Context(),
		metadata,
		conditions,
		pageOptions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrContactsNotFound.WithTrace(err), logger)
		return
	}

	count, err := engine.GetContactsCount(
		c.Request.Context(),
		metadata,
		conditions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotCountContacts.WithTrace(err), logger)
		return
	}

	contracts := common.MapToTypeContracts(contacts, mappings.MapToContactContract)

	response := response.PageModel[response.Contact]{
		Content: contracts,
		Page:    common.GetPageDescriptionFromSearchParams(pageOptions, count),
	}

	c.JSON(http.StatusOK, response)
}

// contactsUpdate will update contact with the given id
// Update contact FullName or Metadata godoc
// @Summary		Update contact FullName or Metadata
// @Description	Update contact FullName or Metadata
// @Tags		Admin
// @Produce		json
// @Param		id path string false "Contact id"
// @Param		UpdateContact body UpdateContact false "FullName and metadata to update"
// @Success		200 {object} response.Contact "Updated contact"
// @Failure		400	"Bad request - Error while parsing UpdateContact from request body or getting id from path"
// @Failure		404	"Not found - Error while getting contact by id"
// @Failure		422	"Unprocessable entity - Incorrect status of contact"
// @Failure 	500	"Internal server error - Error while updating contact"
// @Router	 	/api/v1/admin/contacts/{id} [put]
// @Security	x-auth-xpub
func contactsUpdate(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams UpdateContact
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.WithTrace(err), logger)
		return
	}

	id := c.Param("id")

	contact, err := reqctx.Engine(c).UpdateContact(
		c.Request.Context(),
		id,
		reqParams.FullName,
		&reqParams.Metadata,
	)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrUpdateContact.WithTrace(err), logger)
		return
	}

	contract := mappings.MapToContactContract(contact)

	c.JSON(http.StatusOK, contract)
}

// contactsDelete will delete contact with the given id
// Delete contact godoc
// @Summary		Delete contact
// @Description	Delete contact
// @Tags		Admin
// @Produce		json
// @Param		id path string false "Contact id"
// @Success		200
// @Failure		400	"Bad request - Error while parsing UpdateContact from request body or getting id from path"
// @Failure		404	"Not found - Error while getting contact by id"
// @Failure		422	"Unprocessable entity - Incorrect status of contact"
// @Failure 	500	"Internal server error - Error while updating contact"
// @Failure 	500	"Internal server error - Error while updating contact"
// @Router		/api/v1/admin/contacts/{id} [delete]
// @Security	x-auth-xpub
func contactsDelete(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	id := c.Param("id")

	err := reqctx.Engine(c).DeleteContactByID(
		c.Request.Context(),
		id,
	)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrDeleteContact.WithTrace(err), logger)
		return
	}

	c.Status(http.StatusOK)
}

// contactsReject will reject contact with the given id
// Reject contact with given id godoc
// @Summary		Reject contact
// @Description	Reject contact
// @Tags		Admin
// @Produce		json
// @Param		id path string false "Contact id"
// @Success		200
// @Failure		400	"Bad request - Error while getting id from path"
// @Failure		404	"Not found - Error while getting contact by id"
// @Failure		422	"Unprocessable entity - Incorrect status of contact"
// @Failure 	500	"Internal server error - Error while updating contact"
// @Failure 	500	"Internal server error - Error while changing contact status"
// @Router		/api/v1/admin/invitations/{id} [delete]
// @Security	x-auth-xpub
func contactsReject(c *gin.Context, _ *reqctx.AdminContext) {
	id := c.Param("id")

	_, err := reqctx.Engine(c).AdminChangeContactStatus(
		c.Request.Context(),
		id,
		engine.ContactRejected,
	)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrRejectContact.WithTrace(err), reqctx.Logger(c))
		return
	}

	c.Status(http.StatusOK)
}

// contactsAccept will perform Accept action on contact with the given id
// Perform accept action on contact godoc
// @Summary		Accept contact
// @Description Accept contact
// @Tags		Admin
// @Produce		json
// @Param		id path string false "Contact id"
// @Success		200 {object} response.Contact "Changed contact"
// @Failure		400	"Bad request - Error while getting id from path"
// @Failure		404	"Not found - Error while getting contact by id"
// @Failure		422	"Unprocessable entity - Incorrect status of contact"
// @Failure 	500	"Internal server error - Error while updating contact"
// @Failure 	500	"Internal server error - Error while changing contact status"
// @Router		/api/v1/admin/contact/invitations/{id} [post]
// @Security	x-auth-xpub
func contactsAccept(c *gin.Context, _ *reqctx.AdminContext) {
	id := c.Param("id")

	contact, err := reqctx.Engine(c).AdminChangeContactStatus(
		c.Request.Context(),
		id,
		engine.ContactNotConfirmed,
	)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrAcceptContact.WithTrace(err), reqctx.Logger(c))
		return
	}

	contract := mappings.MapToContactContract(contact)

	c.JSON(http.StatusOK, contract)
}

// contactsCreate will perform create contact action for the given paymail
// @Summary		Create contact
// @Description Create contact
// @Tags		Admin
// @Produce		json
// @Param		paymail path string false "Contact paymail"
// @Success		200 {object} response.Contact "Changed contact"
// @Failure		400	"Bad request - Error while getting paymail from path"
// @Failure		404	"Not found - Error while getting contact requester by paymail"
// @Failure		409	"Contact already exists -  Unable to add duplicate contact"
// @Failure 	500	"Internal server error - Error while adding new contact"
// @Router		/api/v1/admin/contacts/{paymail} [post]
// @Security	x-auth-xpub
func contactsCreate(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	contactPaymail := c.Param("paymail")
	if contactPaymail == "" {
		spverrors.ErrorResponse(c, spverrors.ErrMissingContactPaymailParam, logger)
		return
	}

	var req *CreateContact
	if err := c.Bind(&req); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.WithTrace(err), logger)
		return
	}

	metadata := mappings.MapToMetadata(req.Metadata)
	contact, err := reqctx.Engine(c).AdminCreateContact(c, contactPaymail, req.CreatorPaymail, req.FullName, metadata)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contract := mappings.MapToContactContract(contact)
	c.JSON(http.StatusOK, contract)
}

// contactsConfirm will perform Confirm action on contacts with the given xpub ids and paymails
// Perform confirm action on contacts godoc
// @Summary		Confirm contacts pair
// @Description Marks the contact entries as mutually confirmed, after ensuring the validity of the contact information for both parties.
// @Tags		Admin
// @Produce		json
// @Param		models.AdminConfirmContactPair body models.AdminConfirmContactPair true "Contacts data"
// @Success		200
// @Failure		400	"Bad request - Error while getting data from request body"
// @Failure		404	"Not found - Error, contacts not found"
// @Failure 	500	"Internal server error - Error, confirming contact failed"
// @Router		/api/v1/admin/contacts/confirmations [post]
// @Security	x-auth-xpub
func contactsConfirm(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)

	var reqParams *models.AdminConfirmContactPair
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.WithTrace(err), logger)
		return
	}

	err := reqctx.Engine(c).AdminConfirmContacts(
		c.Request.Context(),
		reqParams.PaymailA,
		reqParams.PaymailB,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.Status(http.StatusOK)
}
