package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux-server/mappings/common"
)

// MapToTransactionContract will map the model from bux to the bux-models contract
func MapToTransactionContract(t *bux.Transaction) *buxmodels.Transaction {
	if t == nil {
		return nil
	}

	model := buxmodels.Transaction{
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
		Status:               string(t.Status),
		TransactionDirection: string(t.Direction),
	}

	processMetadata(t, t.XPubID, &model)

	return &model
}

func processMetadata(t *bux.Transaction, xpubID string, model *buxmodels.Transaction) {
	if len(t.XpubMetadata) > 0 && len(t.XpubMetadata[xpubID]) > 0 {
		if t.Model.Metadata == nil {
			model.Model.Metadata = make(buxmodels.Metadata)
		}
		for key, value := range t.XpubMetadata[xpubID] {
			model.Model.Metadata[key] = value
		}
	}

	model.OutputValue = int64(0)
	if len(t.XpubOutputValue) > 0 && t.XpubOutputValue[xpubID] != 0 {
		model.OutputValue = t.XpubOutputValue[xpubID]
	}

	if model.OutputValue > 0 {
		model.TransactionDirection = "incoming"
	} else {
		model.TransactionDirection = "outgoing"
	}
}

// MapToTransactionBux will map the model from bux-models to the bux contract
func MapToTransactionBux(t *buxmodels.Transaction) *bux.Transaction {
	if t == nil {
		return nil
	}

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

	return &bux.TransactionConfig{
		ChangeDestinations:         mapToBuxDestinations(tx),
		ChangeDestinationsStrategy: bux.ChangeStrategy(tx.ChangeStrategy),
		ChangeMinimumSatoshis:      tx.ChangeMinimumSatoshis,
		ChangeNumberOfDestinations: tx.ChangeNumberOfDestinations,
		ChangeSatoshis:             tx.ChangeSatoshis,
		ExpiresIn:                  tx.ExpiresIn,
		Fee:                        tx.Fee,
		FeeUnit:                    MapToFeeUnitBux(tx.FeeUnit),
		FromUtxos:                  mapToBuxFromUtxos(tx),
		IncludeUtxos:               mapToBuxIncludeUtxos(tx),
		Inputs:                     mapToBuxInputs(tx),
		Outputs:                    mapToBuxOutputs(tx),
		SendAllTo:                  MapToTransactionOutputBux(tx.SendAllTo),
		Sync:                       MapToSyncConfigBux(tx.Sync),
	}
}

func mapToBuxOutputs(tx *buxmodels.TransactionConfig) []*bux.TransactionOutput {
	if tx.Outputs == nil {
		return nil
	}

	outputs := make([]*bux.TransactionOutput, 0)
	for _, output := range tx.Outputs {
		outputs = append(outputs, MapToTransactionOutputBux(output))
	}
	return outputs
}

func mapToBuxInputs(tx *buxmodels.TransactionConfig) []*bux.TransactionInput {
	if tx.Inputs == nil {
		return nil
	}

	inputs := make([]*bux.TransactionInput, 0)
	for _, input := range tx.Inputs {
		inputs = append(inputs, MapToTransactionInputBux(input))
	}
	return inputs
}

func mapToBuxIncludeUtxos(tx *buxmodels.TransactionConfig) []*bux.UtxoPointer {
	if tx.IncludeUtxos == nil {
		return nil
	}

	includeUtxos := make([]*bux.UtxoPointer, 0)
	for _, utxo := range tx.IncludeUtxos {
		includeUtxos = append(includeUtxos, MapToUtxoPointerBux(utxo))
	}
	return includeUtxos
}

func mapToBuxFromUtxos(tx *buxmodels.TransactionConfig) []*bux.UtxoPointer {
	if tx.FromUtxos == nil {
		return nil
	}

	fromUtxos := make([]*bux.UtxoPointer, 0)
	for _, utxo := range tx.FromUtxos {
		fromUtxos = append(fromUtxos, MapToUtxoPointerBux(utxo))
	}
	return fromUtxos
}

func mapToBuxDestinations(tx *buxmodels.TransactionConfig) []*bux.Destination {
	if tx.ChangeDestinations == nil {
		return nil
	}

	destinations := make([]*bux.Destination, 0)
	for _, destination := range tx.ChangeDestinations {
		destinations = append(destinations, MapToDestinationBux(destination))
	}
	return destinations
}

// MapToTransactionConfigContract will map the transaction-config model from bux-models to the bux contract
func MapToTransactionConfigContract(tx *bux.TransactionConfig) *buxmodels.TransactionConfig {
	if tx == nil {
		return nil
	}

	return &buxmodels.TransactionConfig{
		ChangeDestinations:         mapToContractDestinations(tx),
		ChangeStrategy:             string(tx.ChangeDestinationsStrategy),
		ChangeMinimumSatoshis:      tx.ChangeMinimumSatoshis,
		ChangeNumberOfDestinations: tx.ChangeNumberOfDestinations,
		ChangeSatoshis:             tx.ChangeSatoshis,
		ExpiresIn:                  tx.ExpiresIn,
		FeeUnit:                    MapToFeeUnitContract(tx.FeeUnit),
		FromUtxos:                  mapToContractFromUtxos(tx),
		IncludeUtxos:               mapToContractIncludeUtxos(tx),
		Inputs:                     mapToContractInputs(tx),
		Outputs:                    mapToContractOutputs(tx),
		SendAllTo:                  MapToTransactionOutputContract(tx.SendAllTo),
		Sync:                       MapToSyncConfigContract(tx.Sync),
	}
}

func mapToContractOutputs(tx *bux.TransactionConfig) []*buxmodels.TransactionOutput {
	if tx.Outputs == nil {
		return nil
	}

	outputs := make([]*buxmodels.TransactionOutput, 0)
	for _, output := range tx.Outputs {
		outputs = append(outputs, MapToTransactionOutputContract(output))
	}
	return outputs
}

func mapToContractInputs(tx *bux.TransactionConfig) []*buxmodels.TransactionInput {
	if tx.Inputs == nil {
		return nil
	}

	inputs := make([]*buxmodels.TransactionInput, 0)
	for _, input := range tx.Inputs {
		inputs = append(inputs, MapToTransactionInputContract(input))
	}
	return inputs
}

func mapToContractIncludeUtxos(tx *bux.TransactionConfig) []*buxmodels.UtxoPointer {
	if tx.IncludeUtxos == nil {
		return nil
	}

	includeUtxos := make([]*buxmodels.UtxoPointer, 0)
	for _, utxo := range tx.IncludeUtxos {
		includeUtxos = append(includeUtxos, MapToUtxoPointer(utxo))
	}
	return includeUtxos
}

func mapToContractFromUtxos(tx *bux.TransactionConfig) []*buxmodels.UtxoPointer {
	if tx.FromUtxos == nil {
		return nil
	}

	fromUtxos := make([]*buxmodels.UtxoPointer, 0)
	for _, utxo := range tx.FromUtxos {
		fromUtxos = append(fromUtxos, MapToUtxoPointer(utxo))
	}
	return fromUtxos
}

func mapToContractDestinations(tx *bux.TransactionConfig) []*buxmodels.Destination {
	if tx.ChangeDestinations == nil {
		return nil
	}

	destinations := make([]*buxmodels.Destination, 0)
	for _, destination := range tx.ChangeDestinations {
		destinations = append(destinations, MapToDestinationContract(destination))
	}
	return destinations
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
