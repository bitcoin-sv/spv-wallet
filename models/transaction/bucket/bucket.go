package bucket

// Name represents the UTXO bucket where the output belongs to.
type Name string

const (
	// Data represents the bucket for the data only outputs.
	Data Name = "data"
	// BSV represents the bucket for the BSV outputs.
	BSV Name = "bsv"
)
