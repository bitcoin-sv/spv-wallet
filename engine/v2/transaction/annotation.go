package transaction

import (
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

// Annotations represents a transaction metadata that will be used by server to properly handle given transaction.
type Annotations struct {
	Inputs  InputAnnotations
	Outputs OutputsAnnotations
}

// InputAnnotations represents the metadata for chosen inputs. The key is the index of the input.
type InputAnnotations map[int]*InputAnnotation

// InputAnnotation represents the metadata for the input.
type InputAnnotation struct {
	// CustomInstructions has instructions about how to unlock this input.
	CustomInstructions bsv.CustomInstructions
}

// OutputsAnnotations represents the metadata for chosen outputs. The key is the index of the output.
type OutputsAnnotations map[int]*OutputAnnotation

// OutputAnnotation represents the metadata for the output.
type OutputAnnotation struct {
	// Bucket says what type of bucket should this output be stored in.
	Bucket bucket.Name
	// Paymail is available if the output is the paymail output.
	Paymail *PaymailAnnotation
	// CustomInstructions has instructions about how to unlock this output.
	CustomInstructions *bsv.CustomInstructions
}

// PaymailAnnotation is the metadata for the paymail output.
type PaymailAnnotation transaction.PaymailAnnotation

// NewDataOutputAnnotation constructs a new OutputAnnotation for the data output.
func NewDataOutputAnnotation() *OutputAnnotation {
	return &OutputAnnotation{
		Bucket: bucket.Data,
	}
}
