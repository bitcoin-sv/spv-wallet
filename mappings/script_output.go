package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
)

// MapToScriptOutputContract will map the script-output model from bux to the bux-models contract
func MapToScriptOutputContract(so *bux.ScriptOutput) (sc *buxmodels.ScriptOutput) {
	if so == nil {
		return nil
	}

	return &buxmodels.ScriptOutput{
		Address:    so.Address,
		Satoshis:   so.Satoshis,
		Script:     so.Script,
		ScriptType: so.ScriptType,
	}
}

// MapToScriptOutputBux will map the script-output model from bux-models to the bux contract
func MapToScriptOutputBux(so *buxmodels.ScriptOutput) (sc *bux.ScriptOutput) {
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
