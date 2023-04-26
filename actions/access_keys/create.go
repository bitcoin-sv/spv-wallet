package accesskeys

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// create will make a new model using the services defined in the action object
// Create access key godoc
// @Summary		Create access key
// @Description	Create access key
// @Tags		Access-key
// @Produce		json
// @Param		metadata query string false "metadata"
// @Success		201
// @Router		/v1/access-key [post]
// @Security	bux-auth-xpub
func (a *Action) create(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPub, _ := bux.GetXpubFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)

	// params
	metadata := params.GetJSON("metadata")

	// Create a new accessKey
	accessKey, err := a.Services.Bux.NewAccessKey(
		req.Context(),
		reqXPub,
		bux.WithMetadatas(metadata),
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, bux.DisplayModels(accessKey))
}
