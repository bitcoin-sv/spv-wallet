package outlines

import (
	"github.com/bitcoin-sv/go-sdk/script"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
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

func (a annotatedOutputs) totalSatoshis() bsv.Satoshis {
	var total bsv.Satoshis
	for _, out := range a {
		total += bsv.Satoshis(out.Satoshis)
	}
	return total
}

func appendOutput(a annotatedOutputs, satoshis bsv.Satoshis, lockingScript *script.Script, customInstructions bsv.CustomInstructions) annotatedOutputs {
	return append(a, &annotatedOutput{
		OutputAnnotation: &transaction.OutputAnnotation{
			Bucket:             bucket.BSV,
			CustomInstructions: customInstructions,
		},
		TransactionOutput: &sdk.TransactionOutput{
			LockingScript: lockingScript,
			Satoshis:      uint64(satoshis),
		},
	})
}
