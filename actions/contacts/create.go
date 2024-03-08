package contacts

import (
	"errors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// create will make a new model using the services defined in the action object
// Create contact godoc
// @Summary		Create contact
// @Description	Create contact
// @Tags		Contact
// @Produce		json
// @Success		201
// @Router		/v1/contact [post]
// @Security	x-auth-xpub
func (a *Action) create(c *gin.Context) {
	requesterPubKey := c.GetString(auth.ParamXPubKey)

	var req CreateContact
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := req.validate(); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	contact, err := a.Services.SpvWalletEngine.AddContact(
		c.Request.Context(),
		req.FullName, req.Paymail,
		requesterPubKey, req.RequesterFullName, req.RequesterPaymail,
		engine.WithMetadatas(req.Metadata))

	if err != nil && !errors.Is(err, engine.ErrAddingContactRequest) {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	contract := &models.CreateContactResponse{
		Contact: mappings.MapToContactContract(contact),
	}

	if err != nil {
		ai := make(map[string]string)
		ai["warning"] = err.Error()

		contract.AdditionalInfo = ai
	}

	c.JSON(http.StatusCreated, contract)
}
