package outputs

import (
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/evaluation"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
)

// Specifications are representing a client specification for outputs part of the transaction.
type Specifications struct {
	Outputs []Spec
}

// Spec is a specification for a single output of the transaction.
type Spec interface {
	evaluate(ctx evaluation.Context) (annotatedOutputs, error)
}

// NewSpecifications constructs a new Specifications instance with provided outputs specifications.
func NewSpecifications(outputs ...Spec) *Specifications {
	return &Specifications{
		Outputs: outputs,
	}
}

// Add a new output specification to the list of outputs.
func (s *Specifications) Add(output Spec) {
	s.Outputs = append(s.Outputs, output)
}

// Evaluate the outputs specifications and return the transaction outputs and their annotations.
func (s *Specifications) Evaluate(ctx evaluation.Context) ([]*sdk.TransactionOutput, transaction.OutputsAnnotations, error) {
	if s.Outputs == nil {
		return nil, nil, txerrors.ErrDraftRequiresAtLeastOneOutput
	}
	outputs, err := s.evaluate(ctx)
	if err != nil {
		return nil, nil, err
	}

	txOutputs, annotations := outputs.splitIntoTransactionOutputsAndAnnotations()
	return txOutputs, annotations, nil
}

func (s *Specifications) evaluate(ctx evaluation.Context) (annotatedOutputs, error) {
	outputs := make(annotatedOutputs, 0)
	for _, spec := range s.Outputs {
		outs, err := spec.evaluate(ctx)
		if err != nil {
			return nil, spverrors.Wrapf(err, "failed to evaluate output specification %T", spec)
		}
		outputs = append(outputs, outs...)
	}
	return outputs, nil
}

type annotatedOutput struct {
	*transaction.OutputAnnotation
	*sdk.TransactionOutput
}

type annotatedOutputs []*annotatedOutput

func singleAnnotatedOutput(txOut *sdk.TransactionOutput, out *transaction.OutputAnnotation) annotatedOutputs {
	return annotatedOutputs{
		&annotatedOutput{
			OutputAnnotation:  out,
			TransactionOutput: txOut,
		},
	}
}

func (a annotatedOutputs) splitIntoTransactionOutputsAndAnnotations() ([]*sdk.TransactionOutput, transaction.OutputsAnnotations) {
	outputs := make([]*sdk.TransactionOutput, len(a))
	annotationByOutputIndex := make(transaction.OutputsAnnotations)
	for outputIndex, out := range a {
		outputs[outputIndex] = out.TransactionOutput
		if out.OutputAnnotation != nil {
			annotationByOutputIndex[outputIndex] = out.OutputAnnotation
		}
	}
	return outputs, annotationByOutputIndex
}
