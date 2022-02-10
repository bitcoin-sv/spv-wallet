package actions

import (
	"net/http"

	"github.com/BuxOrg/bux-server/dictionary"
	"github.com/julienschmidt/httprouter"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// Health basic request to return a health response
func Health(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	txn := newrelic.FromContext(req.Context())
	txn.Ignore()
	w.WriteHeader(http.StatusOK)
}

// Head is a basic response for any generic HEAD request
func Head(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	txn := newrelic.FromContext(req.Context())
	txn.Ignore()
	w.WriteHeader(http.StatusOK)
}

// NotFound handles all 404 requests
func NotFound(w http.ResponseWriter, req *http.Request) {
	txn := newrelic.FromContext(req.Context())
	txn.Ignore()
	req = newrelic.RequestWithTransactionContext(req, txn)
	ReturnErrorResponse(
		w, req,
		dictionary.GetError(dictionary.ErrorRequestNotFound, req.RequestURI),
		req.RequestURI,
	)
}

// MethodNotAllowed handles all 405 requests
func MethodNotAllowed(w http.ResponseWriter, req *http.Request) {
	txn := newrelic.FromContext(req.Context())
	txn.Ignore()
	req = newrelic.RequestWithTransactionContext(req, txn)
	ReturnErrorResponse(
		w, req,
		dictionary.GetError(dictionary.ErrorMethodNotAllowed, req.Method, req.RequestURI),
		req.Method,
	)
}
