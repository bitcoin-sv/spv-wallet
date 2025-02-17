package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/transactions/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/bsv"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// CreateTransactionOutline creates a transaction outline
func (s *APITransactions) CreateTransactionOutline(c *gin.Context, params api.CreateTransactionOutlineParams) {
	format, err := getOutlineTransactionFormat(params)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	var requestBody api.RequestsTransactionSpecification
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
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), s.logger)
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

	res, err := mapping.TransactionOutlineToResponse(txOutline)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.JSON(200, res)
}

func getOutlineTransactionFormat(params api.CreateTransactionOutlineParams) (bsv.TxHexFormat, error) {
	if params.Format == nil {
		return bsv.TxHexFormatBEEF, nil
	}

	formatValue := *params.Format
	queryFormat := string(formatValue)

	format, err := bsv.ParseTxHexFormat(queryFormat)
	if err == nil {
		return format, spverrors.Wrapf(err, "invalid transaction format [%s] provided for transaction outline", queryFormat)
	}
	return format, nil
}
