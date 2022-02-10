package transactions

import (
	"encoding/json"
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// new will make a new model
// todo: possible duplicate of record.go
func (a *Action) newTransaction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Parse the params
	params := apirouter.GetParams(req)

	// Get the xPub from the request (via authentication)
	reqXPub, _ := bux.GetXpubFromRequest(req)
	xPub, err := a.Services.Bux.GetXpub(req.Context(), reqXPub)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	} else if xPub == nil {
		apirouter.ReturnResponse(w, req, http.StatusForbidden, actions.ErrXpubNotFound.Error())
		return
	}

	// Read transaction config from request body
	// TODO: Austin's params package probably has a better way to do this than
	// marshal/unmarshal... couldn't figure it out
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

	txConfig := bux.TransactionConfig{}
	if err = json.Unmarshal(configBytes, &txConfig); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusBadRequest, actions.ErrBadTxConfig.Error())
		return
	}

	// Record a new transaction (get the hex from parameters)
	var transaction *bux.DraftTransaction
	if transaction, err = a.Services.Bux.NewTransaction(
		req.Context(),
		xPub.RawXpub(),
		&txConfig,
		params.GetJSON(bux.ModelMetadata.String()), // todo: why is this not a point? see other uses of metadata
		// todo: also why is metadata a field? Should use WithMetadata()
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, bux.DisplayModels(transaction))
}
