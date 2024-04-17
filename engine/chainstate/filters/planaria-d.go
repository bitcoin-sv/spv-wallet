package filters

import (
	"strings"

	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/libsv/go-bt"
)

// PlanariaDTemplate string template for a D transaction
const PlanariaDTemplate = "006a223139694733575459537362796f7333754a373333794b347a45696f69314665734e55"

// PlanariaDTemplateAlternate alternate string template for a D transaction
const PlanariaDTemplateAlternate = "6a223139694733575459537362796f7333754a373333794b347a45696f69314665734e55"

// PlanariaD processor
func PlanariaD(tx *chainstate.TxInfo) (*bt.Tx, error) {
	// Loop through all of the outputs and check for pubkeyhash output
	for _, out := range tx.Vout {
		// if any output contains a pubkeyhash output, include this tx in the filter
		if strings.HasPrefix(out.ScriptPubKey.Hex, PlanariaDTemplate) || strings.HasPrefix(out.ScriptPubKey.Hex, PlanariaDTemplateAlternate) {
			return bt.NewTxFromString(tx.Hex)
		}
	}
	return nil, nil
}
