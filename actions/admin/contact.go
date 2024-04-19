package admin

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// contactsSearch will fetch a list of contacts filtered by Metadata and ContactFilters
// Search for contacts filtering by metadata and ContactFilters godoc
// @Summary		Search for contacts
// @Description	Search for contacts
// @Tags		Admin
// @Produce		json
// @Param		SearchRequestParameters body actions.SearchRequestParameters false "Supports targeted resource searches with filters for metadata and custom conditions, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Contact "List of contacts"
// @Failure		400	"Bad request - Error while parsing SearchRequestParameters from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @Router		/v1/admin/contact/search [post]
// @Security	x-auth-xpub
func (a *Action) contactsSearch(c *gin.Context) {
	var reqParams SearchTransactions
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Record a new transaction (get the hex from parameters)a
	contacts, err := a.Services.SpvWalletEngine.GetContacts(
		c.Request.Context(),
		reqParams.Metadata,
		reqParams.Conditions.ToDbConditions(),
		reqParams.QueryParams,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contracts := make([]*models.Contact, 0)
	for _, contact := range contacts {
		contracts = append(contracts, mappings.MapToContactContract(contact))
	}

	c.JSON(http.StatusOK, contracts)
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
// @Failure 	500	"Internal server error - Error while updating contact"
// @Router		/v1/admin/contact/{id} [patch]
// @Security	x-auth-xpub
func (a *Action) contactsUpdate(c *gin.Context) {
	var reqParams UpdateContact
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "id is required")
		return
	}

	contact, err := a.Services.SpvWalletEngine.UpdateContact(
		c.Request.Context(),
		id,
		reqParams.FullName,
		&reqParams.Metadata,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contract := mappings.MapToContactContract(contact)

	c.JSON(http.StatusOK, contract)
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
// @Failure 	500	"Internal server error - Error while changing contact status"
// @Router		/v1/admin/contact/rejected/{id} [patch]
// @Security	x-auth-xpub
func (a *Action) contactsReject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "id is required")
		return
	}

	contact, err := a.Services.SpvWalletEngine.AdminChangeContactStatus(
		c.Request.Context(),
		id,
		engine.ContactRejected,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contract := mappings.MapToContactContract(contact)

	c.JSON(http.StatusOK, contract)
}

// contactsReject will change contact status to unconfirmed
// Change contact status to unconfirmed godoc
// @Summary		Change contact status to unconfirmed
// @Description	Change contact status to unconfirmed
// @Tags		Admin
// @Produce		json
// @Param		id path string false "Contact id"
// @Success		200 {object} models.Contact "Changed contact"
// @Failure		400	"Bad request - Error while getting id from path"
// @Failure 	500	"Internal server error - Error while changing contact status"
// @Router		/v1/admin/contact/unconfirmed/{id} [patch]
// @Security	x-auth-xpub
func (a *Action) contactsUnconfirm(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "id is required")
		return
	}

	contact, err := a.Services.SpvWalletEngine.AdminChangeContactStatus(
		c.Request.Context(),
		id,
		engine.ContactNotConfirmed,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contract := mappings.MapToContactContract(contact)

	c.JSON(http.StatusOK, contract)
}
