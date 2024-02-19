package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapToScriptOutputContract will map the script-output model from spv-wallet to the spv-wallet-models contract
func MapToScriptOutputContract(so *engine.ScriptOutput) (sc *models.ScriptOutput) {
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

// MapToScriptOutputSPV will map the script-output model from spv-wallet-models to the spv-wallet contract
func MapToScriptOutputSPV(so *models.ScriptOutput) (sc *engine.ScriptOutput) {
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
