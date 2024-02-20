package transactions

import (
	"encoding/json"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// newTransaction will create a new transaction
// New transaction godoc
// @Summary		New transaction
// @Description	New transaction
// @Tags		Transactions
// @Produce		json
// @Param		config query string true "transaction config"
// @Param		metadata query string false "metadata"
// @Success		201
// @Router		/v1/transaction [post]
// @Security	x-auth-xpub
func (a *Action) newTransaction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)

	// Get the xPub from the request (via authentication)
	reqXPub, _ := engine.GetXpubFromRequest(req)
	xPub, err := a.Services.SpvWalletEngine.GetXpub(req.Context(), reqXPub)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	} else if xPub == nil {
		apirouter.ReturnResponse(w, req, http.StatusForbidden, actions.ErrXpubNotFound.Error())
		return
	}

	// Read transaction config from request body
	// TODO: consider using go params package functions instead of marshal/unmarshal
	configMap, ok := params.GetJSONOk("config")
	if !ok {
		apirouter.ReturnResponse(w, req, http.StatusBadRequest, actions.ErrTxConfigNotFound.Error())
		return
	}

	var configBytes []byte
	if configBytes, err = json.Marshal(configMap); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusBadRequest, actions.ErrBadTxConfig.Error())
		return
	}

	txContract := models.TransactionConfig{}
	if err = json.Unmarshal(configBytes, &txContract); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusBadRequest, actions.ErrBadTxConfig.Error())
		return
	}

	metadata := params.GetJSON(engine.ModelMetadata.String())
	opts := a.Services.SpvWalletEngine.DefaultModelOptions()
	if metadata != nil {
		opts = append(opts, engine.WithMetadatas(metadata))
	}

	txConfig := mappings.MapTransactionConfigEngineToModel(&txContract)

	// Record a new transaction (get the hex from parameters)
	var transaction *engine.DraftTransaction
	if transaction, err = a.Services.SpvWalletEngine.NewTransaction(
		req.Context(),
		xPub.RawXpub(),
		txConfig,
		opts...,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	contract := mappings.MapToDraftTransactionContract(transaction)

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, contract)
}
