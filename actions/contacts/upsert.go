package contacts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// oldUpsert will add a new contact or modify an existing one.
// Upsert contact godoc
// @Summary		Upsert contact - Use (PUT) /v1/contacts/{paymail} instead.
// @Description	This endpoint has been deprecated. Use (PUT) /v1/contacts/{paymail} instead. Add or update contact. When adding a new contact, the system utilizes Paymail's PIKE capability to dispatch an invitation request, asking the counterparty to include the current user in their contacts.
// @Tags		Contact
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact the user wants to add/modify"
// @Param		UpsertContact body contacts.UpsertContact true "Full name and metadata needed to add/modify contact"
// @Success		201
// @DeprecatedRouter  /v1/contact/{paymail} [PUT]
// @Security	x-auth-xpub
func (a *Action) oldUpsert(c *gin.Context) {
	a.upsertHelper(c, true)
}

// upsertContact will add a new contact or modify an existing one.
// @Summary		Upsert contact
// @Description	Add or update contact. When adding a new contact, the system utilizes Paymail's PIKE capability to dispatch an invitation request, asking the counterparty to include the current user in their contacts.
// @Tags		Contacts
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact the user wants to add/modify"
// @Param		UpsertContact body contacts.UpsertContact true "Full name and metadata needed to add/modify contact"
// @Success		201
// @Router		/v1/contacts/{paymail} [PUT]
// @Security	x-auth-xpub
func (a *Action) upsertContact(c *gin.Context) {
	a.upsertHelper(c, false)

}

func (a *Action) upsertHelper(c *gin.Context, snakeCase bool) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	cPaymail := c.Param("paymail")

	var req UpsertContact
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	if err := req.validate(); err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	contact, err := a.Services.SpvWalletEngine.UpsertContact(
		c.Request.Context(),
		req.FullName, cPaymail,
		reqXPubID, req.RequesterPaymail,
		engine.WithMetadatas(req.Metadata))

	if err != nil && !errors.Is(err, spverrors.ErrAddingContactRequest) {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	if snakeCase {
		res := &models.CreateContactResponse{
			Contact: mappings.MapToOldContactContract(contact),
		}
		if err != nil {
			res.AddAdditionalInfo("warning", err.Error())
		}

		c.JSON(http.StatusOK, res)
		return
	}

	res := &response.CreateContactResponse{
		Contact: mappings.MapToContactContract(contact),
	}
	if err != nil {
		res.AddAdditionalInfo("warning", err.Error())
	}

	c.JSON(http.StatusOK, res)
}
