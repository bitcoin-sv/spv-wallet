package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux-server/mappings/common"
)

// MapToTransactionContract will map the model from bux to the bux-models contract
func MapToTransactionContract(t *bux.Transaction) *buxmodels.Transaction {
	return &buxmodels.Transaction{
		Model:                *common.MapToContract(&t.Model),
		ID:                   t.ID,
		Hex:                  t.Hex,
		XpubInIDs:            t.XpubInIDs,
		XpubOutIDs:           t.XpubOutIDs,
		BlockHash:            t.BlockHash,
		BlockHeight:          t.BlockHeight,
		Fee:                  t.Fee,
		NumberOfInputs:       t.NumberOfInputs,
		NumberOfOutputs:      t.NumberOfOutputs,
		DraftID:              t.DraftID,
		TotalValue:           t.TotalValue,
		OutputValue:          t.OutputValue,
		Status:               string(t.Status),
		TransactionDirection: string(t.Direction),
	}
}

// MapToTransactionBux will map the model from bux-models to the bux contract
func MapToTransactionBux(t *buxmodels.Transaction) *bux.Transaction {
	return &bux.Transaction{
		Model:           *common.MapToModel(&t.Model),
		TransactionBase: bux.TransactionBase{ID: t.ID, Hex: t.Hex},
		XpubInIDs:       t.XpubInIDs,
		XpubOutIDs:      t.XpubOutIDs,
		BlockHash:       t.BlockHash,
		BlockHeight:     t.BlockHeight,
		Fee:             t.Fee,
		NumberOfInputs:  t.NumberOfInputs,
		NumberOfOutputs: t.NumberOfOutputs,
		DraftID:         t.DraftID,
		TotalValue:      t.TotalValue,
		OutputValue:     t.OutputValue,
		Status:          bux.SyncStatus(t.Status),
		Direction:       bux.TransactionDirection(t.TransactionDirection),
	}
}

// MapToTransactionConfigBux will map the transaction-config model from bux to the bux-models contract
func MapToTransactionConfigBux(tx *buxmodels.TransactionConfig) *bux.TransactionConfig {
	if tx == nil {
		return nil
	}

	destinations := make([]*bux.Destination, 0)
	for _, destination := range tx.ChangeDestinations {
		destinations = append(destinations, MapToDestinationBux(destination))
	}

	fromUtxos := make([]*bux.UtxoPointer, 0)
	for _, utxo := range tx.FromUtxos {
		fromUtxos = append(fromUtxos, MapToUtxoPointerBux(utxo))
	}

	includeUtxos := make([]*bux.UtxoPointer, 0)
	for _, utxo := range tx.IncludeUtxos {
		includeUtxos = append(includeUtxos, MapToUtxoPointerBux(utxo))
	}

	inputs := make([]*bux.TransactionInput, 0)
	for _, input := range tx.Inputs {
		inputs = append(inputs, MapToTransactionInputBux(input))
	}

	outputs := make([]*bux.TransactionOutput, 0)
	for _, output := range tx.Outputs {
		outputs = append(outputs, MapToTransactionOutputBux(output))
	}

	return &bux.TransactionConfig{
		ChangeDestinations:         destinations,
		ChangeDestinationsStrategy: bux.ChangeStrategy(tx.ChangeStrategy),
		ChangeMinimumSatoshis:      tx.ChangeMinimumSatoshis,
		ChangeNumberOfDestinations: tx.ChangeNumberOfDestinations,
		ChangeSatoshis:             tx.ChangeSatoshis,
		ExpiresIn:                  tx.ExpiresIn,
		Fee:                        tx.Fee,
		FeeUnit:                    MapToFeeUnitBux(tx.FeeUnit),
		FromUtxos:                  fromUtxos,
		IncludeUtxos:               includeUtxos,
		Inputs:                     inputs,
		Outputs:                    outputs,
		SendAllTo:                  MapToTransactionOutputBux(tx.SendAllTo),
		Sync:                       MapToSyncConfigBux(tx.Sync),
	}
}

// MapToTransactionConfigContract will map the transaction-config model from bux-models to the bux contract
func MapToTransactionConfigContract(tx *bux.TransactionConfig) *buxmodels.TransactionConfig {
	if tx == nil {
		return nil
	}

	destinations := make([]*buxmodels.Destination, 0)
	for _, destination := range tx.ChangeDestinations {
		destinations = append(destinations, MapToDestinationContract(destination))
	}

	fromUtxos := make([]*buxmodels.UtxoPointer, 0)
	for _, utxo := range tx.FromUtxos {
		fromUtxos = append(fromUtxos, MapToUtxoPointer(utxo))
	}

	includeUtxos := make([]*buxmodels.UtxoPointer, 0)
	for _, utxo := range tx.IncludeUtxos {
		includeUtxos = append(includeUtxos, MapToUtxoPointer(utxo))
	}

	inputs := make([]*buxmodels.TransactionInput, 0)
	for _, input := range tx.Inputs {
		inputs = append(inputs, MapToTransactionInputContract(input))
	}

	outputs := make([]*buxmodels.TransactionOutput, 0)
	for _, output := range tx.Outputs {
		outputs = append(outputs, MapToTransactionOutputContract(output))
	}

	return &buxmodels.TransactionConfig{
		ChangeDestinations:         destinations,
		ChangeStrategy:             string(tx.ChangeDestinationsStrategy),
		ChangeMinimumSatoshis:      tx.ChangeMinimumSatoshis,
		ChangeNumberOfDestinations: tx.ChangeNumberOfDestinations,
		ChangeSatoshis:             tx.ChangeSatoshis,
		ExpiresIn:                  tx.ExpiresIn,
		FeeUnit:                    MapToFeeUnitContract(tx.FeeUnit),
		FromUtxos:                  fromUtxos,
		IncludeUtxos:               includeUtxos,
		Inputs:                     inputs,
		Outputs:                    outputs,
		SendAllTo:                  MapToTransactionOutputContract(tx.SendAllTo),
		Sync:                       MapToSyncConfigContract(tx.Sync),
	}
}

// MapToDraftTransactionContract will map the transaction-output model from bux to the bux-models contract
func MapToDraftTransactionContract(tx *bux.DraftTransaction) *buxmodels.DraftTransaction {
	if tx == nil {
		return nil
	}

	return &buxmodels.DraftTransaction{
		Model:         *common.MapToContract(&tx.Model),
		ID:            tx.ID,
		Hex:           tx.Hex,
		XpubID:        tx.XpubID,
		ExpiresAt:     tx.ExpiresAt,
		Configuration: *MapToTransactionConfigContract(&tx.Configuration),
	}
}

// MapToTransactionInputContract will map the transaction-output model from bux-models to the bux contract
func MapToTransactionInputContract(inp *bux.TransactionInput) *buxmodels.TransactionInput {
	if inp == nil {
		return nil
	}

	return &buxmodels.TransactionInput{
		Utxo:        *MapToUtxoContract(&inp.Utxo),
		Destination: *MapToDestinationContract(&inp.Destination),
	}
}

// MapToTransactionInputBux will map the transaction-output model from bux to the bux-models contract
func MapToTransactionInputBux(inp *buxmodels.TransactionInput) *bux.TransactionInput {
	if inp == nil {
		return nil
	}

	return &bux.TransactionInput{
		Utxo:        *MapToUtxoBux(&inp.Utxo),
		Destination: *MapToDestinationBux(&inp.Destination),
	}
}

// MapToTransactionOutputContract will map the transaction-output model from bux to the bux-models contract
func MapToTransactionOutputContract(out *bux.TransactionOutput) *buxmodels.TransactionOutput {
	if out == nil {
		return nil
	}

	scriptOutputs := make([]*buxmodels.ScriptOutput, 0)
	for _, scriptOutput := range out.Scripts {
		scriptOutputs = append(scriptOutputs, MapToScriptOutputContract(scriptOutput))
	}

	return &buxmodels.TransactionOutput{
		OpReturn:     MapToOpReturnContract(out.OpReturn),
		PaymailP4:    MapToPaymailP4Contract(out.PaymailP4),
		Satoshis:     out.Satoshis,
		Script:       out.Script,
		Scripts:      scriptOutputs,
		To:           out.To,
		UseForChange: out.UseForChange,
	}
}

// MapToTransactionOutputBux will map the transaction-output model from bux-models to the bux contract
func MapToTransactionOutputBux(out *buxmodels.TransactionOutput) *bux.TransactionOutput {
	if out == nil {
		return nil
	}

	scriptOutputs := make([]*bux.ScriptOutput, 0)
	for _, scriptOutput := range out.Scripts {
		scriptOutputs = append(scriptOutputs, MapToScriptOutputBux(scriptOutput))
	}

	return &bux.TransactionOutput{
		OpReturn:     MapToOpReturnBux(out.OpReturn),
		PaymailP4:    MapToPaymailP4Bux(out.PaymailP4),
		Satoshis:     out.Satoshis,
		Script:       out.Script,
		Scripts:      scriptOutputs,
		To:           out.To,
		UseForChange: out.UseForChange,
	}
}

// MapToMapProtocolContract will map the transaction-output model from bux to the bux-models contract
func MapToMapProtocolContract(mp *bux.MapProtocol) *buxmodels.MapProtocol {
	if mp == nil {
		return nil
	}

	return &buxmodels.MapProtocol{
		App:  mp.App,
		Keys: mp.Keys,
		Type: mp.Type,
	}
}

// MapToMapProtocolBux will map the transaction-output model from bux-models to the bux contract
func MapToMapProtocolBux(mp *buxmodels.MapProtocol) *bux.MapProtocol {
	if mp == nil {
		return nil
	}

	return &bux.MapProtocol{
		App:  mp.App,
		Keys: mp.Keys,
		Type: mp.Type,
	}
}

// MapToOpReturnContract will map the transaction-output model from bux to the bux-models contract
func MapToOpReturnContract(op *bux.OpReturn) *buxmodels.OpReturn {
	if op == nil {
		return nil
	}

	return &buxmodels.OpReturn{
		Hex:         op.Hex,
		HexParts:    op.HexParts,
		Map:         MapToMapProtocolContract(op.Map),
		StringParts: op.StringParts,
	}
}

// MapToOpReturnBux will map the op-return model from bux-models to the bux contract
func MapToOpReturnBux(op *buxmodels.OpReturn) *bux.OpReturn {
	if op == nil {
		return nil
	}

	return &bux.OpReturn{
		Hex:         op.Hex,
		HexParts:    op.HexParts,
		Map:         MapToMapProtocolBux(op.Map),
		StringParts: op.StringParts,
	}
}
