package mapping

import (
	"maps"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/bsv"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/request"
)

// TransactionOutline maps request's AnnotatedTransaction to outlines.Transaction.
func TransactionOutline(req *request.AnnotatedTransaction) *outlines.Transaction {
	return &outlines.Transaction{
		Hex: bsv.TxHex(req.Hex),
		Annotations: transaction.Annotations{
			Outputs: maps.Collect(func(yield func(int, *transaction.OutputAnnotation) bool) {
				if req.Annotations == nil || len(req.Annotations.Outputs) == 0 {
					return
				}
				for index, output := range req.Annotations.Outputs {
					var paymail *transaction.PaymailAnnotation
					if output.Paymail != nil {
						paymail = &transaction.PaymailAnnotation{
							Receiver:  output.Paymail.Receiver,
							Reference: output.Paymail.Reference,
							Sender:    output.Paymail.Sender,
						}
					}
					yield(index, &transaction.OutputAnnotation{
						Bucket:  output.Bucket,
						Paymail: paymail,
					})
				}
			}),
		},
	}
}
