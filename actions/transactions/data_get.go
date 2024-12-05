package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

func getDataByOutpoint(c *gin.Context, _ *reqctx.UserContext) {
	logger := reqctx.Logger(c)

	var err error
	var outpoint bsv.Outpoint
	outpoint.TxID = c.Param("id")
	outpoint.Vout, err = conv.StringToUint32(c.Param("vout"))
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseParams.Wrap(err), logger)
		return
	}

	// TODO add xpub filtering by userContext.GetXPubID()
	data, err := reqctx.Engine(c).GetTransactionData(
		c,
		outpoint,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}
	if data == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindDataOutpoint, logger)
		return
	}

	c.JSON(http.StatusOK, response.TransactionData{
		Data: string(data),
	})
}
