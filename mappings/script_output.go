package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
)

func MapToScriptOutputContract(so *bux.ScriptOutput) (sc *buxmodels.ScriptOutput) {
	return &buxmodels.ScriptOutput{
		Address:    so.Address,
		Satoshis:   so.Satoshis,
		Script:     so.Script,
		ScriptType: so.ScriptType,
	}
}

func MapToScriptOutputBux(so *buxmodels.ScriptOutput) (sc *bux.ScriptOutput) {
	return &bux.ScriptOutput{
		Address:    so.Address,
		Satoshis:   so.Satoshis,
		Script:     so.Script,
		ScriptType: so.ScriptType,
	}
}
