package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/transactions/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/bsv"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// PostApiV2TransactionsOutlines creates a transaction outline
func (s *APITransactions) PostApiV2TransactionsOutlines(c *gin.Context) {
	format, err := getOutlineTransactionFormat(c)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	var requestBody request.TransactionSpecification
	err = c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), s.logger)
		return
	}

	userCtx := reqctx.GetUserContext(c)
	userID, err := userCtx.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	spec, err := mapping.TransactionSpecificationRequestToOutline(&requestBody, userID)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	var txOutline *outlines.Transaction
	switch format {
	case bsv.TxHexFormatRAW:
		txOutline, err = s.engine.TransactionOutlinesService().CreateRawTx(c, spec)
	case bsv.TxHexFormatBEEF:
		txOutline, err = s.engine.TransactionOutlinesService().CreateBEEF(c, spec)
	}

	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	res := mapping.TransactionOutlineToResponse(txOutline)
	c.JSON(200, res)
}

func getOutlineTransactionFormat(c *gin.Context) (bsv.TxHexFormat, error) {
	queryFormat, _ := c.GetQuery("format")
	if queryFormat == "" {
		return bsv.TxHexFormatBEEF, nil
	}

	format, err := bsv.ParseTxHexFormat(queryFormat)
	if err == nil {
		return format, spverrors.Wrapf(err, "invalid transaction format [%s] provided for transaction outline", queryFormat)
	}
	return format, nil
}
