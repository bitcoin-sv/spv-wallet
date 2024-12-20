package transaction

import (
	"github.com/bitcoin-sv/spv-wallet/models/transaction"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

// Annotations represents a transaction metadata that will be used by server to properly handle given transaction.
type Annotations struct {
	Outputs OutputsAnnotations
}

// OutputAnnotation represents the metadata for the output.
type OutputAnnotation struct {
	// What type of bucket should this output be stored in.
	Bucket bucket.Name
	// Paymail is available if the output is the paymail output.
	Paymail *PaymailAnnotation
}

// PaymailAnnotation is the metadata for the paymail output.
type PaymailAnnotation transaction.PaymailAnnotation

// OutputsAnnotations represents the metadata for chosen outputs. The key is the index of the output.
type OutputsAnnotations map[int]*OutputAnnotation

// NewDataOutputAnnotation constructs a new OutputAnnotation for the data output.
func NewDataOutputAnnotation() *OutputAnnotation {
	return &OutputAnnotation{
		Bucket: bucket.Data,
	}
}
