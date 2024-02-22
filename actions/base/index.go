package base

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// index basic request to /
func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	apirouter.ReturnResponse(
		w, req,
		http.StatusOK,
		map[string]interface{}{"message": "Welcome to the SPV Wallet ✌(◕‿-)✌"},
	)
}
