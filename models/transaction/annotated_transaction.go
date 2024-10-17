package transaction

import (
	"github.com/bitcoin-sv/spv-wallet/models/optional"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

// AnnotatedTransaction represents a transaction with annotations.
type AnnotatedTransaction struct {
	// BEEF is the transaction hex in BEEF format.
	BEEF string `json:"beef"`
	// Annotations is the metadata for the transaction.
	Annotations *Annotations `json:"annotations"`
}

// Annotations represents a transaction metadata that will be used by server to properly handle given transaction.
type Annotations struct {
	Outputs map[int]*OutputAnnotation `json:"outputs"`
}

// OutputAnnotation represents the metadata for the output.
type OutputAnnotation struct {
	// What type of bucket should this output be stored in.
	Bucket bucket.Name `json:"bucket"`
	// Paymail is available if the output is the paymail output.
	Paymail optional.Param[PaymailAnnotation] `json:"paymail,omitempty"`
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
