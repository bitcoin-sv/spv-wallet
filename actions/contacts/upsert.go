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
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// oldUpsert will add a new contact or modify an existing one.
// Upsert contact godoc
// @Summary		Upsert contact - Use (PUT) /api/v1/contacts/{paymail} instead.
// @Description	This endpoint has been deprecated. Use (PUT) /api/v1/contacts/{paymail} instead.
// @Tags		Contact
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact that the user would like to add/modify"
// @Param		UpsertContact body contacts.UpsertContact true "Full name and metadata needed to add/modify contact"
// @Success		201
// @DeprecatedRouter  /v1/contact/{paymail} [put]
// @Security	x-auth-xpub
func oldUpsert(c *gin.Context, userContext *reqctx.UserContext) {
	upsertHelper(c, true, userContext.GetXPubID())
}

// upsertContact will add a new contact or modify an existing one.
// @Summary		Upsert contact
// @Description	Add or update contact. When adding a new contact, the system utilizes Paymail's PIKE capability to dispatch an invitation request, asking the counterparty to include the current user in their contacts.
// @Tags		Contacts
// @Produce		json
// @Param		paymail path string true "Paymail address of the contact that the user would like to add/modify"
// @Param		UpsertContact body contacts.UpsertContact true "Full name and metadata needed to add/modify contact"
// @Success		201
// @Router		/api/v1/contacts/{paymail} [put]
// @Security	x-auth-xpub
func upsertContact(c *gin.Context, userContext *reqctx.UserContext) {
	upsertHelper(c, false, userContext.GetXPubID())
}

func upsertHelper(c *gin.Context, snakeCase bool, xpubID string) {
	logger := reqctx.Logger(c)
	cPaymail := c.Param("paymail")

	var req UpsertContact
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	if err := req.validate(); err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contact, err := reqctx.Engine(c).UpsertContact(
		c.Request.Context(),
		req.FullName, cPaymail,
		xpubID, req.RequesterPaymail,
		engine.WithMetadatas(req.Metadata))

	if err != nil && !errors.Is(err, spverrors.ErrAddingContactRequest) {
		spverrors.ErrorResponse(c, err, logger)
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
