package contacts

import (
	"errors"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

// delete will delete a contact
// Delete contact godoc
// @Summary		Delete contact
// @Description	Delete contact
// @Tags		Contact
// @Produce		json
// @Param		contactId path string true "Contact ID of the contact the user wants to delete"
// @Success		200
// @Failure		404	"Contact not found"
// @Failure		500	"Internal server error"
// @Router		/v1/contact/{contactId} [DELETE]
// @Security	x-auth-xpub
func (a *Action) delete(c *gin.Context) {
	contactID := c.Param(auth.ContactID)

	err := a.Services.SpvWalletEngine.DeleteContact(c.Request.Context(), contactID)

	if err != nil {
		switch {
		case errors.Is(err, engine.ErrContactNotFound):
			c.JSON(http.StatusNotFound, err.Error())
		default:
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.Status(http.StatusOK)
}
