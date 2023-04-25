package accesskeys

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// create will make a new model using the services defined in the action object
// Revoke access key godoc
// @Summary     	Revoke access key
// @Description 	Revoke access key
// @Tags			access-key
// @Produce     	json
// @Param       	id query string true "id"
// @Success     	201
// @Router      	/v1/access-key [delete]
// @Security 		bux-auth-xpub
func (a *Action) revoke(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	reqXPub, _ := bux.GetXpubFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	id := params.GetString("id")

	if id == "" {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, bux.ErrMissingFieldID)
		return
	}

	// Create a new accessKey
	accessKey, err := a.Services.Bux.RevokeAccessKey(
		req.Context(),
		reqXPub,
		id,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, bux.DisplayModels(accessKey))
}
