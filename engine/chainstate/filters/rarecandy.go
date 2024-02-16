package filters

import (
	"strings"

	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/libsv/go-bt"
)

// RareCandyFrogCartelScriptTemplate string template for a Rare Candy Frog Cartel NTF
const RareCandyFrogCartelScriptTemplate = "a914179b4c7a45646a509473df5a444b6e18b723bd148876"

// RareCandyFrogCartel processor
func RareCandyFrogCartel(tx *chainstate.TxInfo) (*bt.Tx, error) {
	// Loop through all the outputs and check for pubkeyhash output
	for _, out := range tx.Vout {
		// if any output contains a pubkeyhash output, include this tx in the filter
		if strings.HasPrefix(out.ScriptPubKey.Hex, RareCandyFrogCartelScriptTemplate) {
			return bt.NewTxFromString(tx.Hex)
		}
	}
	return nil, nil
}
