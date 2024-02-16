package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// status will return the status of the admin login
// Get stats godoc
// @Summary		Get stats
// @Description	Get stats
// @Tags		Admin
// @Produce		json
// @Success		200
// @Router		/v1/admin/stats [get]
// @Security	spv-wallet-auth-xpub
func (a *Action) stats(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stats, err := a.Services.SPV.GetStats(req.Context())
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	contract := mappings.MapToAdminStatsContract(stats)

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, contract)
}
