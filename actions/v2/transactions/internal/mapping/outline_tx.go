package mapping

import (
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/lox"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
	paymailreq "github.com/bitcoin-sv/spv-wallet/models/request/paymail"
	model "github.com/bitcoin-sv/spv-wallet/models/transaction"
	"github.com/samber/lo"
)

// TransactionSpecificationRequestToOutline converts a transaction outline request model to the engine model.
func TransactionSpecificationRequestToOutline(tx *request.TransactionSpecification, userID string) (*outlines.TransactionSpec, error) {
	catcher := lox.NewErrorCollector()

	return &outlines.TransactionSpec{
		UserID: userID,
		Outputs: outlines.OutputsSpec{
			Outputs: lo.Map(
				tx.Outputs,
				lox.MapAndCollect(catcher, lox.MappingFnWithError(outputSpecFromRequest)),
			),
		},
	}, catcher.Error()
}

// TransactionOutlineToResponse converts a transaction outline to a response model.
func TransactionOutlineToResponse(tx *outlines.Transaction) *model.AnnotatedTransaction {
	return &model.AnnotatedTransaction{
		Hex:    string(tx.Hex),
		Format: tx.Hex.Format(),
		Annotations: &model.Annotations{
			Outputs: lo.
				IfF(
					len(tx.Annotations.Outputs) > 0,
					func() map[int]*model.OutputAnnotation {
						return lo.MapValues(tx.Annotations.Outputs, lox.MappingFn(outlineOutputToResponse))
					},
				).Else(nil),
		},
	}
}

func outlineOutputToResponse(from *transaction.OutputAnnotation) *model.OutputAnnotation {
	return &model.OutputAnnotation{
		Bucket: from.Bucket,
		Paymail: lo.
			IfF(
				from.Paymail != nil,
				func() *model.PaymailAnnotation {
					return &model.PaymailAnnotation{
						Sender:    from.Paymail.Sender,
						Receiver:  from.Paymail.Receiver,
						Reference: from.Paymail.Reference,
					}
				},
			).
			Else(nil),
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
