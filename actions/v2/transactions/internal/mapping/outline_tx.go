package mapping

import (
	"errors"
	"strconv"

	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/lox"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/samber/lo"
)

// TransactionSpecificationRequestToOutline converts a transaction outline request model to the engine model.
func TransactionSpecificationRequestToOutline(tx *api.RequestsTransactionSpecification, userID string) (*outlines.TransactionSpec, error) {
	catcher := lox.NewErrorCollector()

	return &outlines.TransactionSpec{
		UserID: userID,
		Outputs: outlines.OutputsSpec{
			Outputs: lo.Map(
				tx.Outputs,
				lox.MapAndCollect(catcher, outputSpecFromRequest),
			),
		},
	}, catcher.Error()
}

// TransactionOutlineToResponse converts a transaction outline to a response model.
func TransactionOutlineToResponse(tx *outlines.Transaction) (api.ModelsAnnotatedTransactionOutline, error) {
	errorCollector := lox.NewErrorCollector()

	return api.ModelsAnnotatedTransactionOutline{
		Hex:    string(tx.Hex),
		Format: api.ModelsAnnotatedTransactionOutlineFormat(tx.Hex.Format()),
		Annotations: &api.ModelsOutlineAnnotations{
			Inputs: lox.Catch(
				errorCollector,
				func() (map[string]api.ModelsInputAnnotation, error) {
					return lox.MapEntriesOrError(tx.Annotations.Inputs, outlineInputEntryToResponse)
				}),
			Outputs: lo.MapEntries(tx.Annotations.Outputs, outlineOutputEntryToResponse),
		},
	}, errorCollector.Error()
}

func outlineInputEntryToResponse(index int, input *transaction.InputAnnotation) (string, api.ModelsInputAnnotation, error) {
	inputAnnotation, err := outlineInputToResponse(input)
	return strconv.Itoa(index), inputAnnotation, err
}

func outlineInputToResponse(item *transaction.InputAnnotation) (api.ModelsInputAnnotation, error) {
	inputAnnotation := api.ModelsInputAnnotation{
		CustomInstructions: api.ModelsCustomInstructions{},
	}

	customInstructions := lo.Map(item.CustomInstructions, lox.MappingFn(customInstructionsToResponse))

	err := inputAnnotation.CustomInstructions.FromModelsSPVWalletCustomInstructions(customInstructions)
	if err != nil {
		return api.ModelsInputAnnotation{}, spverrors.ErrInternal.Wrap(err)
	}

	return inputAnnotation, nil
}

func customInstructionsToResponse(instruction bsv.CustomInstruction) api.ModelsSPVWalletCustomInstruction {
	return api.ModelsSPVWalletCustomInstruction{
		Type:        instruction.Type,
		Instruction: instruction.Instruction,
	}
}

func outlineOutputEntryToResponse(index int, value *transaction.OutputAnnotation) (string, api.ModelsOutputAnnotation) {
	return strconv.Itoa(index), outlineOutputToResponse(value)
}

func outlineOutputToResponse(from *transaction.OutputAnnotation) api.ModelsOutputAnnotation {
	return api.ModelsOutputAnnotation{
		Bucket: api.ModelsOutputAnnotationBucket(from.Bucket),
		CustomInstructions: lo.IfF(
			from.CustomInstructions != nil,
			func() *api.ModelsSPVWalletCustomInstructions {
				return lo.ToPtr(
					lo.Map(*from.CustomInstructions, lox.MappingFn(customInstructionsToResponse)),
				)
			},
		).Else(nil),
		Paymail: lo.IfF(
			from.Paymail != nil,
			func() *api.ModelsPaymailAnnotationDetails {
				return &api.ModelsPaymailAnnotationDetails{
					Sender:    from.Paymail.Sender,
					Receiver:  from.Paymail.Receiver,
					Reference: from.Paymail.Reference,
				}
			},
		).Else(nil),
	}
}

func outputSpecFromRequest(req api.RequestsTransactionOutlineOutputSpecification) (outlines.OutputSpec, error) {
	outputType, err := req.Discriminator()
	if err != nil {
		return nil, spverrors.ErrCannotBindRequest.Wrap(err)
	}

	switch outputType {
	case "op_return":
		return opReturnSpecFromRequest(req)
	case "paymail":
		return paymailSpecFromRequest(req)
	default:
		return nil, spverrors.ErrCannotBindRequest.Wrap(spverrors.Newf("unsupported output type"))
	}
}

func paymailSpecFromRequest(req api.RequestsTransactionOutlineOutputSpecification) (outlines.OutputSpec, error) {
	specification, err := req.AsRequestsPaymailOutputSpecification()
	if err != nil {
		return nil, spverrors.ErrCannotBindRequest.Wrap(err)
	}

	return &outlines.Paymail{
		To:       specification.To,
		Satoshis: bsv.Satoshis(specification.Satoshis),
		From:     specification.From,
	}, nil
}

func opReturnSpecFromRequest(req api.RequestsTransactionOutlineOutputSpecification) (outlines.OutputSpec, error) {
	specification, err := req.AsRequestsOpReturnOutputSpecification()
	if err != nil {
		return nil, spverrors.ErrCannotBindRequest.Wrap(err)
	}

	var dataType string
	if specification.DataType == nil {
		dataType = ""
	} else {
		dataType = string(*specification.DataType)
	}

	switch dataType {
	case "", "strings":
		v, err := specification.Data.AsRequestsOpReturnStringsOutput()
		if err != nil {
			return nil, spverrors.ErrCannotBindRequest.Wrap(err)
		}
		return &outlines.OpReturn{
			DataType: outlines.DataTypeStrings,
			Data:     v,
		}, nil
	case "hexes":
		v, err := specification.Data.AsRequestsOpReturnHexesOutput()
		if err != nil {
			return nil, spverrors.ErrCannotBindRequest.Wrap(err)
		}

		return &outlines.OpReturn{
			DataType: outlines.DataTypeHexes,
			Data:     v,
		}, nil
	default:
		return nil, errors.New("unsupported output type")
	}
}
