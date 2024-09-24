package merkleroots

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// TODO: handle errors
func get(c *gin.Context, userContext *reqctx.UserContext, appConfig *config.AppConfig) {
	batchSize := c.Query("batchSize")
	lastEvaluatedKey := c.Query("lastEvaluatedKey")
	client := http.Client{}
	bhsUrl, err := createBHSURL(appConfig, "/chain/merkleroot")
	if err != nil {
		c.JSON(500, err)
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
		c.JSON(500, err)
		return
	}

	if appConfig.BHS.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+appConfig.BHS.AuthToken)
	}

	res, err := client.Do(req)
	if err != nil {
		c.JSON(500, err)
		return
	}

	var response any
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, response)
}

func createBHSURL(appConfig *config.AppConfig, endpointPath string) (*url.URL, error) {
	url, err := url.Parse(appConfig.BHS.URL + "/api/" + config.APIVersion + endpointPath)
	if err != nil {
		return nil, err
	}

	return url, nil
}
