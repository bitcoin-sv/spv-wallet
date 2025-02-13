package transaction

import (
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/optional"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

// AnnotatedTransaction represents a transaction with annotations.
type AnnotatedTransaction struct {
	// Hex is the transaction in binary format specified by type.
	Hex string `json:"hex"`
	// Format is the format of the transaction hex ex. BEEF, RAW.
	Format string `json:"format"`
	// Annotations is the metadata for the transaction.
	Annotations *Annotations `json:"annotations"`
}

// Annotations represents a transaction metadata that will be used by server to properly handle given transaction.
type Annotations struct {
	Outputs map[int]*OutputAnnotation `json:"outputs"`
	Inputs  map[int]*InputAnnotation  `json:"inputs"`
}

// OutputAnnotation represents the metadata for the output.
type OutputAnnotation struct {
	// What type of bucket should this output be stored in.
	Bucket bucket.Name `json:"bucket"`
	// Paymail is available if the output is the paymail output.
	Paymail optional.Param[PaymailAnnotation] `json:"paymail,omitempty"`
	// CustomInstructions has instructions about how to unlock this output.
	CustomInstructions optional.Param[bsv.CustomInstructions] `json:"customInstructions,omitempty"`
}

// InputAnnotation represents the metadata for the input.
type InputAnnotation struct {
	// CustomInstructions has instructions about how to unlock this input.
	CustomInstructions bsv.CustomInstructions `json:"customInstructions"`
}

// PaymailAnnotation is the metadata for the paymail output.
type PaymailAnnotation struct {
	// Receiver is the paymail address of the receiver.
	Receiver string `json:"receiver"`
	// Reference is the reference number used for paymail transaction.
	Reference string `json:"reference"`
	// Sender is the paymail address of the sender.
	Sender string `json:"sender"`
}
