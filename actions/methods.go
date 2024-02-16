package actions

import (
	"encoding/json"
	"net/http"

	spvwalletmodels "github.com/bitcoin-sv/bux-models"
	"github.com/bitcoin-sv/spv-wallet/dictionary"
	"github.com/julienschmidt/httprouter"
	"github.com/mrz1836/go-datastore"
	"github.com/mrz1836/go-parameters"
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

// GetQueryParameters get all filtering parameters related to the db query
func GetQueryParameters(params *parameters.Params) (*datastore.QueryParams, *spvwalletmodels.Metadata, *map[string]interface{}, error) {
	var queryParams *datastore.QueryParams
	jsonQueryParams, ok := params.GetJSONOk("params")
	if ok {
		p, err := json.Marshal(jsonQueryParams)
		if err != nil {
			return nil, nil, nil, err
		}
		err = json.Unmarshal(p, &queryParams)
		if err != nil {
			return nil, nil, nil, err
		}
	} else {
		queryParams = &datastore.QueryParams{
			Page:          params.GetInt("page"),
			PageSize:      params.GetInt("page_size"),
			OrderByField:  params.GetString("order_by_field"),
			SortDirection: params.GetString("sort_direction"),
		}
	}

	metadataReq := params.GetJSON(MetadataField)
	var metadata *spvwalletmodels.Metadata
	if len(metadataReq) > 0 {
		// marshal the metadata into the Metadata model
		metaJSON, _ := json.Marshal(metadataReq) //nolint:errchkjson // ignore for now
		if err := json.Unmarshal(metaJSON, &metadata); err != nil {
			return nil, nil, nil, err
		}
	}
	conditionsReq := params.GetJSON("conditions")
	var conditions *map[string]interface{}
	if len(conditionsReq) > 0 {
		// marshal the conditions into the Map
		conditionsJSON, _ := json.Marshal(conditionsReq) //nolint:errchkjson // ignore for now
		if err := json.Unmarshal(conditionsJSON, &conditions); err != nil {
			return nil, nil, nil, err
		}
	}

	return queryParams, metadata, conditions, nil
}
