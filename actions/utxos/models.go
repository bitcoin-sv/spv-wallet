package utxos

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

// CountUtxo is the model containing filters for counting utxos
type CountUtxo struct {
	// Custom conditions used for filtering the search results
	Conditions map[string]interface{} `json:"conditions"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata"`
}
