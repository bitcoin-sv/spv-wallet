package utxos

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

// CountUtxo is the model containing filters for counting utxos
type CountUtxo struct {
	Conditions map[string]interface{} `json:"conditions"`
	Metadata   engine.Metadata        `json:"metadata"`
}
