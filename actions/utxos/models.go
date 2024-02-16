package utxos

import "github.com/BuxOrg/bux"

type CountUtxo struct {
	Conditions map[string]interface{} `json:"conditions"`
	Metadata   bux.Metadata           `json:"metadata"`
}
