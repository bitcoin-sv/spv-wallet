package mapping

import (
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/mapper"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
	paymailreq "github.com/bitcoin-sv/spv-wallet/models/request/paymail"
	model "github.com/bitcoin-sv/spv-wallet/models/transaction"
	"github.com/samber/lo"
)

// TransactionSpecificationRequestToOutline converts a transaction outline request model to the engine model.
func TransactionSpecificationRequestToOutline(tx *request.TransactionSpecification, userID string) *outlines.TransactionSpec {
	return &outlines.TransactionSpec{
		UserID: userID,
		Outputs: outlines.OutputsSpec{
			Outputs: lo.Map(tx.Outputs, mapper.MapWithoutIndex(transactionRequestOutputsToOutline)),
		},
	}
}

// TransactionOutlineToResponse converts a transaction outline to a response model.
func TransactionOutlineToResponse(tx *outlines.Transaction) *model.AnnotatedTransaction {
	var annotations model.Annotations
	if len(tx.Annotations.Outputs) > 0 {
		annotations.Outputs = lo.MapValues(tx.Annotations.Outputs, mapper.MapWithoutIndex(outlineOutputToResponse))
	}

	return &model.AnnotatedTransaction{
		Hex:         string(tx.Hex),
		Format:      tx.Hex.Format(),
		Annotations: &annotations,
	}
}

func transactionRequestOutputsToOutline(val request.Output) outlines.OutputSpec {
	spec, err := outputSpecFromRequest(val)
	// TODO: handle error
	if err != nil {
		return nil
	}
	return spec
}

func outlineOutputToResponse(from *transaction.OutputAnnotation) *model.OutputAnnotation {
	outputAnnotation := &model.OutputAnnotation{
		Bucket: from.Bucket,
	}
	if from.Paymail != nil {
		outputAnnotation.Paymail = &model.PaymailAnnotation{
			Sender:    from.Paymail.Sender,
			Receiver:  from.Paymail.Receiver,
			Reference: from.Paymail.Reference,
		}
	}

	return outputAnnotation
}

func outputSpecFromRequest(req request.Output) (outlines.OutputSpec, error) {
	switch o := req.(type) {
	case opreturn.Output:
		out := outlines.OpReturn(o)
		return &out, nil
	case paymailreq.Output:
		out := outlines.Paymail(o)
		return &out, nil
	default:
		return nil, errors.New("unsupported output type")
	}
}
