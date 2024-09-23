package transaction

// Bucket represents the UTXO bucket where the output belongs to.
type Bucket string

const (
	// BucketData represents the bucket for the data only outputs.
	BucketData Bucket = "data"
	// BucketBSV represents the bucket for the BSV outputs.
	BucketBSV Bucket = "bsv"
)

// Annotations represents a transaction metadata that will be used by server to properly handle given transaction.
type Annotations struct {
	Outputs OutputsAnnotations
}

// OutputAnnotation represents the metadata for the output.
type OutputAnnotation struct {
	// What type of bucket should this output be stored in.
	Bucket Bucket
	// Paymail is available if the output is the paymail output.
	Paymail *PaymailAnnotation
}

// PaymailAnnotation is the metadata for the paymail output.
type PaymailAnnotation struct {
	Receiver  string
	Reference string
}

// OutputsAnnotations represents the metadata for chosen outputs. The key is the index of the output.
type OutputsAnnotations map[int]*OutputAnnotation

// NewDataOutputAnnotation constructs a new OutputAnnotation for the data output.
func NewDataOutputAnnotation() *OutputAnnotation {
	return &OutputAnnotation{
		Bucket: BucketData,
	}
}
