package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/gin-gonic/gin"
)

// contactsSearch will fetch a list of contacts filtered by Metadata and ContactFilters
// Search for contacts filtering by metadata and ContactFilters godoc
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
func (a *Action) contactsSearch(c *gin.Context) {
	var reqParams filter.SearchContacts
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	reqParams.DefaultsIfNil()

	contacts, err := a.Services.SpvWalletEngine.GetContacts(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		conditions,
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	contracts := mappings.MapToContactContracts(contacts)

	count, err := a.Services.SpvWalletEngine.GetContactsCount(
		c.Request.Context(),
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

// contactsUpdate will update contact with the given id
// Update contact FullName or Metadata godoc
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
func (a *Action) contactsUpdate(c *gin.Context) {
	var reqParams UpdateContact
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	id := c.Param("id")

	contact, err := a.Services.SpvWalletEngine.UpdateContact(
		c.Request.Context(),
		id,
		reqParams.FullName,
		&reqParams.Metadata,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
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
// @Router		/v1/admin/contact/{id} [delete]
// @Security	x-auth-xpub
func (a *Action) contactsDelete(c *gin.Context) {
	id := c.Param("id")

	err := a.Services.SpvWalletEngine.DeleteContact(
		c.Request.Context(),
		id,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
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
// @Success		200 {object} models.Contact "Rejected contact"
// @Failure		400	"Bad request - Error while getting id from path"
// @Failure		404	"Not found - Error while getting contact by id"
// @Failure		422	"Unprocessable entity - Incorrect status of contact"
// @Failure 	500	"Internal server error - Error while updating contact"
// @Failure 	500	"Internal server error - Error while changing contact status"
// @Router		/v1/admin/contact/rejected/{id} [patch]
// @Security	x-auth-xpub
func (a *Action) contactsReject(c *gin.Context) {
	id := c.Param("id")

	contact, err := a.Services.SpvWalletEngine.AdminChangeContactStatus(
		c.Request.Context(),
		id,
		engine.ContactRejected,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	contract := mappings.MapToContactContract(contact)

	c.JSON(http.StatusOK, contract)
}

// contactsAccept will perform Accept action on contact with the given id
// Perform accept action on contact godoc
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
func (a *Action) contactsAccept(c *gin.Context) {
	id := c.Param("id")

	contact, err := a.Services.SpvWalletEngine.AdminChangeContactStatus(
		c.Request.Context(),
		id,
		engine.ContactNotConfirmed,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	contract := mappings.MapToContactContract(contact)

	c.JSON(http.StatusOK, contract)
}
