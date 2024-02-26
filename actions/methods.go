package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/dictionary"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/mrz1836/go-datastore"
)

// RequestParameters is a struct for handling basic request parameters
type RequestParameters struct {
	Conditions  map[string]interface{} `json:"conditions"`
	Metadata    engine.Metadata        `json:"metadata"`
	QueryParams datastore.QueryParams  `json:"query_params"`
}

// StatusOK is a basic response which sets the status to 200
func StatusOK(c *gin.Context) {
	c.Status(http.StatusOK)
}

// NotFound handles all 404 requests
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, dictionary.GetError(dictionary.ErrorRequestNotFound, c.Request.RequestURI))
}

// MethodNotAllowed handles all 405 requests
func MethodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, dictionary.GetError(dictionary.ErrorMethodNotAllowed, c.Request.Method, c.Request.RequestURI))
}

// GetQueryParameters get all filtering parameters related to the db query
func GetQueryParameters(c *gin.Context) (*datastore.QueryParams, *engine.Metadata, *map[string]interface{}, error) {
	var requestParameters RequestParameters
	if err := c.Bind(&requestParameters); err != nil {
		return nil, nil, nil, err
	}

	if requestParameters.Conditions == nil {
		requestParameters.Conditions = make(map[string]interface{})
	}

	return &requestParameters.QueryParams, &requestParameters.Metadata, &requestParameters.Conditions, nil
}
