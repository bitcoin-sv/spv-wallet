package filters

import (
	"strings"

	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/libsv/go-bt"
)

// PlanariaBTemplate string template for a B transaction
const PlanariaBTemplate = "006a2231394878696756345179427633744870515663554551797131707a5a56646f417574"

// PlanariaBTemplateAlternate alternate string template for a B transaction
const PlanariaBTemplateAlternate = "6a2231394878696756345179427633744870515663554551797131707a5a56646f417574"

// PlanariaB processor
func PlanariaB(tx *chainstate.TxInfo) (*bt.Tx, error) {
	// Loop through all the outputs and check for pubkeyhash output
	for _, out := range tx.Vout {
		// if any output contains a pubkeyhash output, include this tx in the filter
		if strings.HasPrefix(out.ScriptPubKey.Hex, PlanariaBTemplate) || strings.HasPrefix(out.ScriptPubKey.Hex, PlanariaBTemplateAlternate) {
			return bt.NewTxFromString(tx.Hex)
		}
	}
	return nil, nil
}
