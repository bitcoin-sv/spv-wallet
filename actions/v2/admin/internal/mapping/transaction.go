package mapping

import (
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/bsv"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/bitcoin-sv/spv-wallet/lox"
	bsvmodel "github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
	"github.com/samber/lo"
)

// RecordedOutline maps domain RecordedOutline to api.ModelsRecordedOutline.
func RecordedOutline(r *txmodels.RecordedOutline) api.ModelsRecordedOutline {
	return api.ModelsRecordedOutline{
		TxID: r.TxID,
	}
}

// RequestsTransactionOutlineToOutline maps request's AnnotatedTransaction to outlines.Transaction.
func RequestsTransactionOutlineToOutline(req *api.RequestsRecordTransactionOutlineForUser) (*outlines.Transaction, error) {
	errorCollector := lox.NewErrorCollector()

	return &outlines.Transaction{
		Hex: bsv.TxHex(req.Hex),
		Annotations: transaction.Annotations{
			Outputs: lo.
				IfF(
					req.Annotations != nil,
					lox.CatchFn(errorCollector, func() (transaction.OutputsAnnotations, error) {
						return lox.MapEntriesOrError(req.Annotations.Outputs, mapOutputAnnotationEntry)
					}),
				).Else(nil),
		},
	}, errorCollector.Error()
}

func mapOutputAnnotationEntry(key string, value api.ModelsOutputAnnotation) (uint32, *transaction.OutputAnnotation, error) {
	var vout uint32
	if n, err := fmt.Sscanf(key, "%d", &vout); err != nil {
		return 0, nil, spverrors.ErrCannotMapFromModel.Wrap(err)
	} else if n != 1 {
		return 0, nil, spverrors.ErrCannotMapFromModel.Wrap(spverrors.Newf("failed to parse vout from key %s", key))
	}
	return vout, annotatedOutputToOutline(value), nil
}

func annotatedOutputToOutline(from api.ModelsOutputAnnotation) *transaction.OutputAnnotation {
	return &transaction.OutputAnnotation{
		Bucket: bucket.Name(from.Bucket),
		Paymail: lo.IfF(
			from.Paymail != nil,
			func() *transaction.PaymailAnnotation {
				return &transaction.PaymailAnnotation{
					Sender:    from.Paymail.Sender,
					Receiver:  from.Paymail.Receiver,
					Reference: from.Paymail.Reference,
				}
			},
		).Else(nil),
		CustomInstructions: lo.IfF(
			from.CustomInstructions != nil,
			func() *bsvmodel.CustomInstructions {
				return lo.ToPtr(
					bsvmodel.CustomInstructions(lo.Map(*from.CustomInstructions, lox.MappingFn(requestToCustomResponse))),
				)
			},
		).Else(nil),
	}
}

func requestToCustomResponse(instruction api.ModelsSPVWalletCustomInstruction) bsvmodel.CustomInstruction {
	return bsvmodel.CustomInstruction{
		Type:        instruction.Type,
		Instruction: instruction.Instruction,
	}
}
