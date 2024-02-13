package mappings

import (
	"github.com/BuxOrg/bux"
	spvwalletmodels "github.com/BuxOrg/bux-models"
)

// MapToScriptOutputContract will map the script-output model from spv-wallet to the spv-wallet-models contract
func MapToScriptOutputContract(so *bux.ScriptOutput) (sc *spvwalletmodels.ScriptOutput) {
	if so == nil {
		return nil
	}

	return &spvwalletmodels.ScriptOutput{
		Address:    so.Address,
		Satoshis:   so.Satoshis,
		Script:     so.Script,
		ScriptType: so.ScriptType,
	}
}

// MapToScriptOutputSPV will map the script-output model from spv-wallet-models to the spv-wallet contract
func MapToScriptOutputSPV(so *spvwalletmodels.ScriptOutput) (sc *bux.ScriptOutput) {
	if so == nil {
		return nil
	}

	return &bux.ScriptOutput{
		Address:    so.Address,
		Satoshis:   so.Satoshis,
		Script:     so.Script,
		ScriptType: so.ScriptType,
	}
}
