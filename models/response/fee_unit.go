package response

// FeeUnit is a model that represents a fee unit (simplified version of fee unit from go-bt).
type FeeUnit struct {
	// Satoshis is a fee unit satoshis amount.
	Satoshis int `json:"satoshis" example:"1"`
	// Bytes is a fee unit bytes representation.
	Bytes int `json:"bytes" example:"1000"`
}
