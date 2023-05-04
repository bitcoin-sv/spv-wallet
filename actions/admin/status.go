package admin

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// status will return the status of the admin login
// Get status godoc
// @Summary		Get status
// @Description	Get status
// @Tags		Admin
// @Produce		json
// @Success		200
// @Router		/v1/admin/status [get]
// @Security	bux-auth-xpub
func (a *Action) status(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, true)
}
