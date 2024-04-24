package actions

import (
	"fmt"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/dictionary"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/gin-gonic/gin"
)

// SearchRequestParameters is a struct for handling request parameters for search requests
type SearchRequestParameters struct {
	// Custom conditions used for filtering the search results
	Conditions map[string]interface{} `json:"conditions" swaggertype:"object,string" example:"testColumn:testValue"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
	// Pagination and sorting options to streamline data exploration and analysis
	QueryParams datastore.QueryParams `json:"params" swaggertype:"object,string" example:"page:1,page_size:10,order_by_field:created_at,order_by_direction:desc"`
}

// CountRequestParameters is a struct for handling request parameters for count requests
type CountRequestParameters struct {
	// Custom conditions used for filtering the search results
	Conditions map[string]interface{} `json:"conditions"  swaggertype:"object,string" example:"testColumn:testValue"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
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
		err = fmt.Errorf("error occurred while binding request parameters: %w", err)
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
		err = fmt.Errorf("error occurred while binding request parameters: %w", err)
		return nil, nil, err
	}

	if requestParameters.Conditions == nil {
		requestParameters.Conditions = make(map[string]interface{})
	}

	return &requestParameters.Metadata, &requestParameters.Conditions, nil
}
