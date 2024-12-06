package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// contactsSearchOld will fetch a list of contacts filtered by Metadata and ContactFilters
// Search for contacts filtering by metadata and ContactFilters godoc
// @DeprecatedRouter /v1/admin/contact/search [post]
// @Summary		Search for contacts
// @Description	Search for contacts
// @Tags		Admin
// @Produce		json
// @Param		SearchContacts body filter.SearchContacts false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} models.SearchContactsResponse "List of contacts"
// @Failure		400	"Bad request - Error while parsing SearchContacts from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @Router		/v1/admin/contact/search [post]
// @Security	x-auth-xpub
func contactsSearchOld(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	engine := reqctx.Engine(c)
	var reqParams filter.SearchContacts
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	reqParams.DefaultsIfNil()

	contacts, err := engine.GetContacts(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contracts := mappings.MapToOldContactContracts(contacts)

	count, err := engine.GetContactsCount(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	response := models.SearchContactsResponse{
		Content: contracts,
		Page:    common.GetPageFromQueryParams(reqParams.QueryParams, count),
	}

	c.JSON(http.StatusOK, response)
}

// contactsUpdateOld will update contact with the given id
// Update contact FullName or Metadata godoc
// @DeprecatedRouter /v1/admin/contact/{id} [patch]
// @Summary		Update contact FullName or Metadata
// @Description	Update contact FullName or Metadata
// @Tags		Admin
// @Produce		json
// @Param		id path string false "Contact id"
// @Param		UpdateContact body UpdateContact false "FullName and metadata to update"
// @Success		200 {object} models.Contact "Updated contact"
// @Failure		400	"Bad request - Error while parsing UpdateContact from request body or getting id from path"
// @Failure		404	"Not found - Error while getting contact by id"
// @Failure		422	"Unprocessable entity - Incorrect status of contact"
// @Failure 	500	"Internal server error - Error while updating contact"
// @Router		/v1/admin/contact/{id} [patch]
// @Security	x-auth-xpub
func contactsUpdateOld(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams UpdateContact
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
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
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contract := mappings.MapToContactContract(contact)

	c.JSON(http.StatusOK, contract)
}

// contactsDeleteOld will delete contact with the given id
// Delete contact godoc
// @DeprecatedRouter /v1/admin/contact/{id} [delete]
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
// @Router		/v1/admin/contact/{id} [delete]
// @Security	x-auth-xpub
func contactsDeleteOld(c *gin.Context, _ *reqctx.AdminContext) {
	id := c.Param("id")

	err := reqctx.Engine(c).DeleteContactByID(
		c.Request.Context(),
		id,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.Status(http.StatusOK)
}

// contactsRejectOld will reject contact with the given id
// Reject contact with given id godoc
// @DeprecatedRouter /v1/admin/contact/rejected/{id} [patch]
// @Summary		Reject contact
// @Description	Reject contact
// @Tags		Admin
// @Produce		json
// @Param		id path string false "Contact id"
// @Success		200 {object} models.Contact "Rejected contact"
// @Failure		400	"Bad request - Error while getting id from path"
// @Failure		404	"Not found - Error while getting contact by id"
// @Failure		422	"Unprocessable entity - Incorrect status of contact"
// @Failure 	500	"Internal server error - Error while updating contact"
// @Failure 	500	"Internal server error - Error while changing contact status"
// @Router		/v1/admin/contact/rejected/{id} [patch]
// @Security	x-auth-xpub
func contactsRejectOld(c *gin.Context, _ *reqctx.AdminContext) {
	id := c.Param("id")

	contact, err := reqctx.Engine(c).AdminChangeContactStatus(
		c.Request.Context(),
		id,
		engine.ContactRejected,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	contract := mappings.MapToContactContract(contact)

	c.JSON(http.StatusOK, contract)
}

// contactsAcceptOld will perform Accept action on contact with the given id
// Perform accept action on contact godoc
// @DeprecatedRouter /v1/admin/contact/accepted/{id} [patch]
// @Summary		Accept contact
// @Description Accept contact
// @Tags		Admin
// @Produce		json
// @Param		id path string false "Contact id"
// @Success		200 {object} models.Contact "Changed contact"
// @Failure		400	"Bad request - Error while getting id from path"
// @Failure		404	"Not found - Error while getting contact by id"
// @Failure		422	"Unprocessable entity - Incorrect status of contact"
// @Failure 	500	"Internal server error - Error while updating contact"
// @Failure 	500	"Internal server error - Error while changing contact status"
// @Router		/v1/admin/contact/accepted/{id} [patch]
// @Security	x-auth-xpub
func contactsAcceptOld(c *gin.Context, _ *reqctx.AdminContext) {
	id := c.Param("id")

	contact, err := reqctx.Engine(c).AdminChangeContactStatus(
		c.Request.Context(),
		id,
		engine.ContactNotConfirmed,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	contract := mappings.MapToContactContract(contact)

	c.JSON(http.StatusOK, contract)
}

// contactsConfirmOld will perform Confirm action on contacts with the given xpub ids and paymails
// Perform confirm action on contacts godoc
// @Summary		Confirm contacts
// @Description Confirm contacts
// @Tags		Admin
// @Produce		json
// @Param		[]models.AdminConfirmContactPair body []models.AdminConfirmContactPair true "Contacts data"
// @Success		200
// @Failure		400	"Bad request - Error while getting data from request body"
// @Failure		404	"Not found - Error, contacts not found"
// @Failure 	500	"Internal server error - Error, confirming contact failed"
// @Router		/v1/admin/contacts/confirmations [post]
// @Security	x-auth-xpub
func contactsConfirmOld(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)

	var reqParams *models.AdminConfirmContactPair
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.WithTrace(err), logger)
		return
	}

	contacts := mappings.MapToEngineContractAdminConfirmContactPair(reqParams)

	err := reqctx.Engine(c).AdminConfirmContacts(
		c.Request.Context(),
		contacts,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.Status(http.StatusOK)
}
