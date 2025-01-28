package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	model "github.com/bitcoin-sv/spv-wallet/models/transaction"
	"github.com/samber/lo"
)

func AnnotatedTransactionToOutline(tx *model.AnnotatedTransaction) *outlines.Transaction {
	return &outlines.Transaction{
		BEEF: tx.BEEF,
		Annotations: transaction.Annotations{
			Outputs: lo.MapValues(tx.Annotations.Outputs, annotatedOutputToOutline),
		},
	}
}

func annotatedOutputToOutline(from *model.OutputAnnotation, _ int) *transaction.OutputAnnotation {
	return &transaction.OutputAnnotation{
		Bucket: from.Bucket,
		Paymail: &transaction.PaymailAnnotation{
			Sender:    from.Paymail.Sender,
			Receiver:  from.Paymail.Receiver,
			Reference: from.Paymail.Reference,
		},
	}
}
