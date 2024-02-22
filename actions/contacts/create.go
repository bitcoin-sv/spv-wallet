package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
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
func (c *Action) create(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	params := apirouter.GetParams(req)

	fullName := params.GetString("full_name")
	paymail := params.GetString("paymail")
	pubKey := params.GetString("pubKey")
	metadata := params.GetJSON("metadata")

	opts := c.Services.SpvWalletEngine.DefaultModelOptions()
	if metadata != nil {
		opts = append(opts, engine.WithMetadatas(metadata))
	}

	contact, err := c.Services.SpvWalletEngine.NewContact(req.Context(), fullName, paymail, pubKey, opts...)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	contract := mappings.MapToContactContract(contact)

	apirouter.ReturnResponse(w, req, http.StatusCreated, contract)
}
