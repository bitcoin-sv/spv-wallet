package mapping

import (
	"strconv"

	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/bsv"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/lox"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
	"github.com/samber/lo"
)

// AnnotatedTransactionRequestToOutline maps request's AnnotatedTransaction to outlines.Transaction.
func AnnotatedTransactionRequestToOutline(req *api.ApiComponentsRequestsAnnotatedTransaction) (*outlines.Transaction, error) {
	errorCollector := lox.NewErrorCollector()

	return &outlines.Transaction{
		Hex: bsv.TxHex(req.Hex),
		Annotations: transaction.Annotations{
			Outputs: lo.
				IfF(
					req.Annotations != nil,
					lox.CatchFn(errorCollector, func() (transaction.OutputsAnnotations, error) {
						return lox.MapEntriesOrError(*req.Annotations.Outputs, mapOutputAnnotationEntry)
					}),
				).Else(nil),
		},
	}, errorCollector.Error()
}

func mapOutputAnnotationEntry(key string, value api.ApiComponentsModelsOutputAnnotation) (int, *transaction.OutputAnnotation, error) {
	index, err := strconv.Atoi(key)
	if err != nil {
		return 0, nil, spverrors.ErrCannotMapFromModel.Wrap(err)
	}
	return index, annotatedOutputToOutline(value), nil
}

func annotatedOutputToOutline(from api.ApiComponentsModelsOutputAnnotation) *transaction.OutputAnnotation {
	return &transaction.OutputAnnotation{
		Bucket: bucket.Name(from.Bucket),
		Paymail: lo.
			IfF(
				from.Paymail != nil,
				func() *transaction.PaymailAnnotation {
					return &transaction.PaymailAnnotation{
						Sender:    from.Paymail.Sender,
						Receiver:  from.Paymail.Receiver,
						Reference: from.Paymail.Reference,
					}
				},
			).Else(nil),
	}
}
