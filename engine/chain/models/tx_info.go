package chainmodels

// TXInfo is the struct that represents the transaction information from ARC
type TXInfo struct {
	BlockHash   string   `json:"blockHash,omitempty"`
	BlockHeight int64    `json:"blockHeight,omitempty"`
	ExtraInfo   string   `json:"extraInfo,omitempty"`
	MerklePath  string   `json:"merklePath,omitempty"`
	Timestamp   string   `json:"timestamp,omitempty"`
	TXStatus    TXStatus `json:"txStatus,omitempty"`
	TxID        string   `json:"txid,omitempty"`
}

// NotFound presents a convention to indicate that the transaction is not found
func (t *TXInfo) NotFound() bool {
	return t == nil
}
