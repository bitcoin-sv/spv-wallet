package contacts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// upsert will add a new contact or modify an existing one.
// Upsert contact godoc
// @Summary		Upsert contact
// @Description	Add or update contact. For new contact send request to add current user as contact
// @Tags		Contact
// @Produce		json
// @Param		UpsertContact body contacts.UpsertContact true "Full name and metadata needed to add/modify contact"
// @Success		201
// @Router		/v1/contact/{paymail} [PUT]
// @Security	x-auth-xpub
func (a *Action) upsert(c *gin.Context) {
	requesterPubKey := c.GetString(auth.ParamXPubKey)
	cPaymail := c.GetString("paymail")

	var req UpsertContact
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := req.validate(); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	contact, err := a.Services.SpvWalletEngine.UpsertContact(
		c.Request.Context(),
		req.FullName, cPaymail,
		requesterPubKey,
		engine.WithMetadatas(req.Metadata))

	if err != nil && !errors.Is(err, engine.ErrAddingContactRequest) {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := &models.CreateContactResponse{
		Contact: mappings.MapToContactContract(contact),
	}

	if err != nil {
		response.AddAdditionalInfo("warning", err.Error())
	}

	c.JSON(http.StatusOK, response)
}
