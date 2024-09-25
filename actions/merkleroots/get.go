package merkleroots

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

func get(c *gin.Context, userContext *reqctx.UserContext, appConfig *config.AppConfig) {
	client := &http.Client{}
	logger := reqctx.Logger(c)
	bhsOK := CheckBlockHeaderServiceStatus(c, appConfig.BHS, client, logger)
	if !bhsOK {
		spverrors.ErrorResponse(c, spverrors.ErrBHSUnreachable, logger)
		return
	}

	batchSize := c.Query("batchSize")
	lastEvaluatedKey := c.Query("lastEvaluatedKey")
	bhsUrl, err := createBHSURL(appConfig, "/chain/merkleroot")
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	query := url.Values{}
	if batchSize != "" {
		query.Add("batchSize", batchSize)
	}
	if lastEvaluatedKey != "" {
		query.Add("lastEvaluatedKey", lastEvaluatedKey)
	}

	bhsUrl.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(c, http.MethodGet, bhsUrl.String(), nil)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	if appConfig.BHS.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+appConfig.BHS.AuthToken)
	}

	res, err := client.Do(req)
	if res != nil {
		defer func() {
			_ = res.Body.Close()
		}()
	}
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	if res.StatusCode != http.StatusOK {
		mapBHSErrorResponseToSpverror(c, res, logger)
		return
	}

	var response any
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrBHSParsingResponse, logger)
		return
	}

	c.JSON(http.StatusOK, response)
}

// createBHSURL parses Block Header Url from configuration and constructs a valid
// endpoint with provided endpointPath variable
func createBHSURL(appConfig *config.AppConfig, endpointPath string) (*url.URL, error) {
	url, err := url.Parse(appConfig.BHS.URL + "/api/" + config.APIVersion + endpointPath)
	if err != nil {
		return nil, spverrors.ErrBHSBadURL
	}

	return url, nil
}
