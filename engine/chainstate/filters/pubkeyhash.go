package filters

import (
	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/libsv/go-bt"
)

// PubKeyHash processor
func PubKeyHash(tx *chainstate.TxInfo) (*bt.Tx, error) {
	// log.Printf("Attempting to filter for pubkeyhash: %#v", tx)
	// Loop through all the outputs and check for pubkeyhash output
	for _, out := range tx.Vout {
		// if any output contains a pubkeyhash output, include this tx in the filter
		if out.ScriptPubKey.Type == "pubkeyhash" {
			return bt.NewTxFromString(tx.Hex)
		}
	}
	return nil, nil
}
