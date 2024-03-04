package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/gin-gonic/gin"
)

// create will make a new model using the services defined in the action object
// Create contact godoc
// @Summary		Create contact
// @Description	Create contact
// @Tags		Contact
// @Produce		json
// @Param		fullName query string false "fullName"
// @Param		paymail query string false "paymail"
// @Param		pubKey query string false "pubKey"
// @Param		metadata query string false "metadata"
// @Success		201
// @Router		/v1/contact [post]
// @Security	bux-auth-xpub
func (a *Action) create(c *gin.Context) {
	requesterPubKey := c.GetString("pubKey")

	var req CreateContact
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	contact, err := a.Services.SpvWalletEngine.AddContact(
		c.Request.Context(),
		req.FullName, req.Paymail,
		requesterPubKey, req.RequesterFullName, req.RequesterPaymail,
		engine.WithMetadatas(req.Metadata))

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	contract := mappings.MapToContactContract(contact)
	c.JSON(http.StatusCreated, contract)
}
