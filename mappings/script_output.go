package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// MapToScriptOutputContract will map the script-output model from spv-wallet to the spv-wallet-models contract
func MapToScriptOutputContract(so *engine.ScriptOutput) (sc *response.ScriptOutput) {
	if so == nil {
		return nil
	}

	return &response.ScriptOutput{
		Address:    so.Address,
		Satoshis:   so.Satoshis,
		Script:     so.Script,
		ScriptType: so.ScriptType,
	}
}

// MapScriptOutputModelToEngine will map the script-output model from spv-wallet-models to the spv-wallet contract
func MapScriptOutputModelToEngine(so *response.ScriptOutput) (sc *engine.ScriptOutput) {
	if so == nil {
		return nil
	}

	return &engine.ScriptOutput{
		Address:    so.Address,
		Satoshis:   so.Satoshis,
		Script:     so.Script,
		ScriptType: so.ScriptType,
	}
}
