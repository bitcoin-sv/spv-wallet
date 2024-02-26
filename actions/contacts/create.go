package contacts

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
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

	fullName := c.GetString("full_name")
	paymail := c.GetString("paymail")
	pubKey := c.GetString("pubKey")

	var requestBody CreateContact

	opts := a.Services.SpvWalletEngine.DefaultModelOptions()
	if requestBody.Metadata != nil {
		opts = append(opts, engine.WithMetadatas(requestBody.Metadata))
	}

	contact, err := a.Services.SpvWalletEngine.NewContact(c.Request.Context(), fullName, paymail, pubKey, opts...)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	contract := mappings.MapToContactContract(contact)

	c.JSON(http.StatusCreated, contract)
}
