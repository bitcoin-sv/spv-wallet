package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapToOldScriptOutputContract will map the script-output model from spv-wallet to the spv-wallet-models contract
func MapToOldScriptOutputContract(so *engine.ScriptOutput) (sc *models.ScriptOutput) {
	if so == nil {
		return nil
	}

	return &models.ScriptOutput{
		Address:    so.Address,
		Satoshis:   so.Satoshis,
		Script:     so.Script,
		ScriptType: so.ScriptType,
	}
}

// MapOldScriptOutputModelToEngine will map the script-output model from spv-wallet-models to the spv-wallet contract
func MapOldScriptOutputModelToEngine(so *models.ScriptOutput) (sc *engine.ScriptOutput) {
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
