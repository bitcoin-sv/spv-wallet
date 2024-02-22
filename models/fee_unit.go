package models

// FeeUnit is a model that represents a fee unit (simplified version of fee unit from go-bt).
type FeeUnit struct {
	// Satoshis is a fee unit satoshis amount.
	Satoshis int `json:"satoshis"`
	// Bytes is a fee unit bytes representation.
	Bytes int `json:"bytes"`
}
