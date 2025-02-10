package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/bsv"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
	"github.com/samber/lo"
	"strconv"
)

// AnnotatedTransactionRequestToOutline maps request's AnnotatedTransaction to outlines.Transaction.
func AnnotatedTransactionRequestToOutline(req *api.ApiComponentsRequestsAnnotatedTransaction) *outlines.Transaction {
	return &outlines.Transaction{
		Hex: bsv.TxHex(req.Hex),
		Annotations: transaction.Annotations{
			Outputs: lo.
				IfF(
					req.Annotations != nil,
					func() transaction.OutputsAnnotations {
						return lo.MapEntries(*req.Annotations.Outputs, annotatedOutputToOutline)
					},
				).Else(nil),
		},
	}
}

func annotatedOutputToOutline(key string, from api.ApiComponentsModelsOutputAnnotation) (int, *transaction.OutputAnnotation) {
	// TODO: Errorcollector
	intKey, _ := strconv.Atoi(key)
	return intKey, &transaction.OutputAnnotation{
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
