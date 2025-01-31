package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/engine/v2/bsv"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/lox"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	model "github.com/bitcoin-sv/spv-wallet/models/transaction"
	"github.com/samber/lo"
)

// AnnotatedTransactionRequestToOutline maps request's AnnotatedTransaction to outlines.Transaction.
func AnnotatedTransactionRequestToOutline(req *request.AnnotatedTransaction) *outlines.Transaction {
	return &outlines.Transaction{
		Hex: bsv.TxHex(req.Hex),
		Annotations: transaction.Annotations{
			Outputs: lo.
				IfF(
					req.Annotations != nil && len(req.Annotations.Outputs) > 0,
					func() transaction.OutputsAnnotations {
						return lo.MapValues(req.Annotations.Outputs, lox.MappingFn(annotatedOutputToOutline))
					},
				).Else(nil),
		},
	}
}

func annotatedOutputToOutline(from *model.OutputAnnotation) *transaction.OutputAnnotation {
	return &transaction.OutputAnnotation{
		Bucket: from.Bucket,
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
