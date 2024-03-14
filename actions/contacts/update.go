package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// update will update an existing model
// Update Contact godoc
// @Summary		Update contact
// @Description	Update contact
// @Tags		Contacts
// @Produce		json
// @Param		metadata body string true "Contacts Metadata"
// @Success		200
// @Router		/v1/contact [patch]
// @Security	x-auth-xpub
func (a *Action) update(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var requestBody UpdateContact

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if requestBody.XPubID == "" {
		c.JSON(http.StatusBadRequest, "Id is missing")
	}

	contact, err := a.Services.SpvWalletEngine.UpdateContact(c.Request.Context(), requestBody.FullName, requestBody.PubKey, reqXPubID, requestBody.Paymail, engine.ContactStatus(requestBody.Status), engine.WithMetadatas(requestBody.Metadata))

	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	contract := mappings.MapToContactContract(contact)

	c.JSON(http.StatusOK, contract)

}
