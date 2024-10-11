package chainmodels

import "time"

// TXInfo is the struct that represents the transaction information from ARC
type TXInfo struct {
	BlockHash   string    `json:"blockHash,omitempty"`
	BlockHeight int64     `json:"blockHeight,omitempty"`
	ExtraInfo   string    `json:"extraInfo,omitempty"`
	MerklePath  string    `json:"merklePath,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
	TXStatus    TXStatus  `json:"txStatus,omitempty"`
	TxID        string    `json:"txid,omitempty"`
}

// Found presents a convention to indicate that the transaction is known by ARC
func (t *TXInfo) Found() bool {
	return t != nil
}

func (t *TXInfo) IsSuccess() bool {
	return t.BlockHeight > 0
}
