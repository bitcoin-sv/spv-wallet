package mapping

import (
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
	paymailreq "github.com/bitcoin-sv/spv-wallet/models/request/paymail"
	model "github.com/bitcoin-sv/spv-wallet/models/transaction"
	"github.com/samber/lo"
)

// Request is a transaction outline request model.
type Request request.TransactionSpecification

// ToEngine converts a transaction outline request model to the engine model.
func TransactionRequestToOutline(tx *request.TransactionSpecification, userID string) *outlines.TransactionSpec {
	return &outlines.TransactionSpec{
		UserID: userID,
		Outputs: outlines.OutputsSpec{
			Outputs: lo.Map(tx.Outputs, transactionRequestOutputsToOutline),
		},
	}
}

func transactionRequestOutputsToOutline(val request.Output, _ int) outlines.OutputSpec {
	spec, err := outputSpecFromRequest(val)
	// TODO: handle error
	if err != nil {
		return nil
	}
	return spec
}

// ToResponse converts a transaction outline to a response model.
func TransactionOutlineToResponse(tx *outlines.Transaction) *model.AnnotatedTransaction {
	return &model.AnnotatedTransaction{
		BEEF: tx.BEEF,
		Annotations: &model.Annotations{
			Outputs: lo.MapValues(tx.Annotations.Outputs, outlineOutputToResponse),
		},
	}
}

func outlineOutputToResponse(from *transaction.OutputAnnotation, _ int) *model.OutputAnnotation {
	return &model.OutputAnnotation{
		Bucket: from.Bucket,
		Paymail: &model.PaymailAnnotation{
			Receiver:  from.Paymail.Receiver,
			Reference: from.Paymail.Reference,
			Sender:    from.Paymail.Sender,
		},
	}
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
