package annotatedtx

import (
	"maps"

	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
	model "github.com/bitcoin-sv/spv-wallet/models/transaction"
)

// Request represents a request model for recording a transaction outline.
type Request model.AnnotatedTransaction

// ToEngine converts a request model to the engine model.
func (req Request) ToEngine() *outlines.Transaction {
	return &outlines.Transaction{
		BEEF: req.BEEF,
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
