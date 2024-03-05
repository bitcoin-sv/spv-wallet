package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
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

	var contact engine.Contact

	err := c.ShouldBindJSON(&contact)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	pubKey := c.GetString(auth.ParamXPubKey)

	var requestBody CreateContact

	newContact, err := a.Services.SpvWalletEngine.NewContact(c.Request.Context(), contact.FullName, contact.Paymail, pubKey, engine.WithMetadatas(requestBody.Metadata))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	contract := mappings.MapToContactContract(newContact)

	c.JSON(http.StatusCreated, contract)
}
