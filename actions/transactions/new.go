package transactions

import (
	"encoding/json"
	"net/http"

	"github.com/BuxOrg/bux"
	spvwalletmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/spv-wallet/actions"
	"github.com/BuxOrg/spv-wallet/mappings"
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
// @Security	spv-wallet-auth-xpub
func (a *Action) newTransaction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)

	// Get the xPub from the request (via authentication)
	reqXPub, _ := bux.GetXpubFromRequest(req)
	xPub, err := a.Services.SPV.GetXpub(req.Context(), reqXPub)
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

	txContract := spvwalletmodels.TransactionConfig{}
	if err = json.Unmarshal(configBytes, &txContract); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusBadRequest, actions.ErrBadTxConfig.Error())
		return
	}

	metadata := params.GetJSON(bux.ModelMetadata.String())
	opts := a.Services.SPV.DefaultModelOptions()
	if metadata != nil {
		opts = append(opts, bux.WithMetadatas(metadata))
	}

	txConfig := mappings.MapToTransactionConfigSPV(&txContract)

	// Record a new transaction (get the hex from parameters)
	var transaction *bux.DraftTransaction
	if transaction, err = a.Services.SPV.NewTransaction(
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
