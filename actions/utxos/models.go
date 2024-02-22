package utxos

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

type CountUtxo struct {
	Conditions map[string]interface{} `json:"conditions"`
	Metadata   engine.Metadata        `json:"metadata"`
}
