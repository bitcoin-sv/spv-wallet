package actions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/dictionary"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/gin-gonic/gin"
	"github.com/mrz1836/go-datastore"
)

// SearchRequestParameters is a struct for handling request parameters for search requests
type SearchRequestParameters struct {
	Conditions  map[string]interface{} `json:"conditions"`
	Metadata    engine.Metadata        `json:"metadata"`
	QueryParams datastore.QueryParams  `json:"params"`
}

// CountRequestParameters is a struct for handling request parameters for count requests
type CountRequestParameters struct {
	Metadata   engine.Metadata        `json:"metadata"`
	Conditions map[string]interface{} `json:"conditions"`
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

// GetSearchQueryParameters get all filtering parameters related to the db query
func GetSearchQueryParameters(c *gin.Context) (*datastore.QueryParams, *engine.Metadata, *map[string]interface{}, error) {
	var requestParameters SearchRequestParameters
	if err := c.Bind(&requestParameters); err != nil {
		return nil, nil, nil, err
	}

	if requestParameters.Conditions == nil {
		requestParameters.Conditions = make(map[string]interface{})
	}

	return &requestParameters.QueryParams, &requestParameters.Metadata, &requestParameters.Conditions, nil
}

// GetCountQueryParameters get all filtering parameters related to the db query
func GetCountQueryParameters(c *gin.Context) (*engine.Metadata, *map[string]interface{}, error) {
	var requestParameters CountRequestParameters
	if err := c.Bind(&requestParameters); err != nil {
		return nil, nil, err
	}

	if requestParameters.Conditions == nil {
		requestParameters.Conditions = make(map[string]interface{})
	}

	return &requestParameters.Metadata, &requestParameters.Conditions, nil
}
