package accesskeys

import (
	"encoding/json"
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux/utils"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// list will fetch a list of access keys
func (a *Action) list(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	reqXPub, _ := bux.GetXpubFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	metadataReq := params.GetJSON(bux.ModelMetadata.String())
	var metadata *bux.Metadata
	if len(metadataReq) > 0 {
		// marshal the metadata into the Metadata model
		metaJSON, _ := json.Marshal(metadataReq) // nolint: errchkjson // ignore for now
		_ = json.Unmarshal(metaJSON, &metadata)
	}

	// Record a new transaction (get the hex from parameters)a
	var err error
	var accessKeys []*bux.AccessKey
	if accessKeys, err = a.Services.Bux.GetAccessKeys(
		req.Context(),
		utils.Hash(reqXPub),
		metadata,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, bux.DisplayModels(accessKeys))
}
