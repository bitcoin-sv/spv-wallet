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
	var annotations transaction.Annotations
	if req.Annotations != nil && len(req.Annotations.Outputs) > 0 {
		annotations.Outputs = lo.MapValues(req.Annotations.Outputs, lox.MappingFn(annotatedOutputToOutline))
	}

	return &outlines.Transaction{
		Hex:         bsv.TxHex(req.Hex),
		Annotations: annotations,
	}
}

// AnnotatedTransactionToOutline maps AnnotatedTransaction model to Transaction engine model
func AnnotatedTransactionToOutline(tx *model.AnnotatedTransaction) *outlines.Transaction {
	var annotations transaction.Annotations
	if len(tx.Annotations.Outputs) > 0 {
		annotations.Outputs = lo.MapValues(tx.Annotations.Outputs, lox.MappingFn(annotatedOutputToOutline))
	}

	return &outlines.Transaction{
		Hex:         bsv.TxHex(tx.Hex),
		Annotations: annotations,
	}
}

func annotatedOutputToOutline(from *model.OutputAnnotation) *transaction.OutputAnnotation {
	outputAnnotation := &transaction.OutputAnnotation{
		Bucket: from.Bucket,
	}
	if from.Paymail != nil {
		outputAnnotation.Paymail = &transaction.PaymailAnnotation{
			Sender:    from.Paymail.Sender,
			Receiver:  from.Paymail.Receiver,
			Reference: from.Paymail.Reference,
		}
	}

	return outputAnnotation
}
