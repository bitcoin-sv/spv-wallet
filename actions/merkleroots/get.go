package merkleroots

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func get(c *gin.Context, userContext *reqctx.UserContext) {
	client := resty.New()
	logger := reqctx.Logger(c)
	appConfig := reqctx.AppConfig(c)

	res, err := getMerkleRootsFromBHS(client, appConfig, logger, c.Request.URL.Query())

	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.JSON(http.StatusOK, res)
}
