package outlines

import (
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
)

// OutputsSpec are representing a client specification for outputs part of the transaction.
type OutputsSpec struct {
	Outputs []OutputSpec
}

// OutputSpec is a specification for a single output of the transaction.
type OutputSpec interface {
	evaluate(ctx *evaluationContext) (annotatedOutputs, error)
}

// NewOutputsSpecs constructs a new OutputsSpec instance with provided outputs specifications.
func NewOutputsSpecs(outputs ...OutputSpec) OutputsSpec {
	return OutputsSpec{
		Outputs: outputs,
	}
}

// Add a new output specification to the list of outputs.
func (s *OutputsSpec) Add(output OutputSpec) {
	s.Outputs = append(s.Outputs, output)
}

// Evaluate the outputs specifications and return the transaction outputs and their annotations.
func (s *OutputsSpec) Evaluate(ctx *evaluationContext) ([]*sdk.TransactionOutput, transaction.OutputsAnnotations, error) {
	if s.Outputs == nil {
		return nil, nil, txerrors.ErrTxOutlineRequiresAtLeastOneOutput
	}
	outputs, err := s.evaluate(ctx)
	if err != nil {
		return nil, nil, err
	}

	txOutputs, annotations := outputs.splitIntoTransactionOutputsAndAnnotations()
	return txOutputs, annotations, nil
}

func (s *OutputsSpec) evaluate(ctx *evaluationContext) (annotatedOutputs, error) {
	if len(s.Outputs) == 0 {
		return nil, txerrors.ErrTxOutlineRequiresAtLeastOneOutput
	}

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

type annotatedOutputs []*annotatedOutput

type annotatedOutput struct {
	*transaction.OutputAnnotation
	*sdk.TransactionOutput
}

func singleAnnotatedOutput(txOut *sdk.TransactionOutput, out *transaction.OutputAnnotation) annotatedOutputs {
	return annotatedOutputs{
		&annotatedOutput{
			OutputAnnotation:  out,
			TransactionOutput: txOut,
		},
	}
}

func (a annotatedOutputs) splitIntoTransactionOutputsAndAnnotations() ([]*sdk.TransactionOutput, transaction.OutputsAnnotations) {
	return a.toTransactionOutputs(), a.toAnnotations()
}

func (a annotatedOutputs) toTransactionOutputs() []*sdk.TransactionOutput {
	outputs := make([]*sdk.TransactionOutput, len(a))
	for i, out := range a {
		outputs[i] = out.TransactionOutput
	}
	return outputs
}

func (a annotatedOutputs) toAnnotations() transaction.OutputsAnnotations {
	annotations := make(transaction.OutputsAnnotations)
	for outputIndex, out := range a {
		if out.OutputAnnotation != nil {
			annotations[outputIndex] = out.OutputAnnotation
		}
	}
	return annotations
}
