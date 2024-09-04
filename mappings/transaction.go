package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// MapToTransactionContract will map the model from spv-wallet to the spv-wallet-models contract
func MapToTransactionContract(t *engine.Transaction) *response.Transaction {
	if t == nil {
		return nil
	}

	model := response.Transaction{
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
		Status:               string(t.TxStatus),
		TransactionDirection: string(t.Direction),
	}

	processMetadata(t, t.XPubID, &model)
	processOutputValue(t, t.XPubID, &model)

	return &model
}

// MapToTransactionContractForAdmin will map the model from spv-wallet to the spv-wallet-response contract for admin
func MapToTransactionContractForAdmin(t *engine.Transaction) *response.Transaction {
	if t == nil {
		return nil
	}

	model := response.Transaction{
		Model:           *common.MapToContract(&t.Model),
		ID:              t.ID,
		Hex:             t.Hex,
		XpubInIDs:       t.XpubInIDs,
		XpubOutIDs:      t.XpubOutIDs,
		BlockHash:       t.BlockHash,
		BlockHeight:     t.BlockHeight,
		Fee:             t.Fee,
		NumberOfInputs:  t.NumberOfInputs,
		NumberOfOutputs: t.NumberOfOutputs,
		DraftID:         t.DraftID,
		TotalValue:      t.TotalValue,
		Status:          string(t.Status),
		Outputs:         t.XpubOutputValue,
	}

	processMetadata(t, t.XPubID, &model)

	return &model
}

func processMetadata(t *engine.Transaction, xpubID string, model *response.Transaction) {
	if len(t.XpubMetadata) > 0 && len(t.XpubMetadata[xpubID]) > 0 {
		if t.Model.Metadata == nil {
			model.Model.Metadata = make(models.Metadata)
		}
		for key, value := range t.XpubMetadata[xpubID] {
			model.Model.Metadata[key] = value
		}
	}
}

func processOutputValue(t *engine.Transaction, xpubID string, model *response.Transaction) {
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

// MapTransactionModelToEngine will map the model from spv-wallet-response to the spv-wallet contract
func MapTransactionModelToEngine(t *response.Transaction) *engine.Transaction {
	if t == nil {
		return nil
	}

	return &engine.Transaction{
		Model:           *common.MapToModel(&t.Model),
		TransactionBase: engine.TransactionBase{ID: t.ID, Hex: t.Hex},
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
		Status:          engine.SyncStatus(t.Status),
		Direction:       engine.TransactionDirection(t.TransactionDirection),
	}
}

// MapTransactionConfigEngineToModel will map the transaction-config model from spv-wallet to the spv-wallet-response contract
func MapTransactionConfigEngineToModel(tx *response.TransactionConfig) *engine.TransactionConfig {
	if tx == nil {
		return nil
	}

	return &engine.TransactionConfig{
		ChangeDestinations:         mapToEngineDestinations(tx),
		ChangeDestinationsStrategy: engine.ChangeStrategy(tx.ChangeStrategy),
		ChangeMinimumSatoshis:      tx.ChangeMinimumSatoshis,
		ChangeNumberOfDestinations: tx.ChangeNumberOfDestinations,
		ChangeSatoshis:             tx.ChangeSatoshis,
		ExpiresIn:                  tx.ExpiresIn,
		Fee:                        tx.Fee,
		FeeUnit:                    MapFeeUnitModelToEngine(tx.FeeUnit),
		FromUtxos:                  mapToEngineFromUtxos(tx),
		IncludeUtxos:               mapIncludeUtxosModelToEngine(tx),
		Inputs:                     mapToEngineInputs(tx),
		Outputs:                    mapToEngineOutputs(tx),
		SendAllTo:                  MapTransactionOutputModelToEngine(tx.SendAllTo),
		Sync:                       MapSyncConfigModelToEngine(tx.Sync),
	}
}

func mapToEngineOutputs(tx *response.TransactionConfig) []*engine.TransactionOutput {
	if tx.Outputs == nil {
		return nil
	}

	outputs := make([]*engine.TransactionOutput, 0)
	for _, output := range tx.Outputs {
		outputs = append(outputs, MapTransactionOutputModelToEngine(output))
	}
	return outputs
}

func mapToEngineInputs(tx *response.TransactionConfig) []*engine.TransactionInput {
	if tx.Inputs == nil {
		return nil
	}

	inputs := make([]*engine.TransactionInput, 0)
	for _, input := range tx.Inputs {
		inputs = append(inputs, MapTransactionInputModelToEngine(input))
	}
	return inputs
}

func mapIncludeUtxosModelToEngine(tx *response.TransactionConfig) []*engine.UtxoPointer {
	if tx.IncludeUtxos == nil {
		return nil
	}

	includeUtxos := make([]*engine.UtxoPointer, 0)
	for _, utxo := range tx.IncludeUtxos {
		includeUtxos = append(includeUtxos, MapUtxoPointerModelToEngine(utxo))
	}
	return includeUtxos
}

func mapToEngineFromUtxos(tx *response.TransactionConfig) []*engine.UtxoPointer {
	if tx.FromUtxos == nil {
		return nil
	}

	fromUtxos := make([]*engine.UtxoPointer, 0)
	for _, utxo := range tx.FromUtxos {
		fromUtxos = append(fromUtxos, MapUtxoPointerModelToEngine(utxo))
	}
	return fromUtxos
}

func mapToEngineDestinations(tx *response.TransactionConfig) []*engine.Destination {
	if tx.ChangeDestinations == nil {
		return nil
	}

	destinations := make([]*engine.Destination, 0)
	for _, destination := range tx.ChangeDestinations {
		destinations = append(destinations, MapDestinationModelToEngine(destination))
	}
	return destinations
}

// MapToTransactionConfigContract will map the transaction-config model from spv-wallet-models to the spv-wallet contract
func MapToTransactionConfigContract(tx *engine.TransactionConfig) *response.TransactionConfig {
	if tx == nil {
		return nil
	}

	return &response.TransactionConfig{
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

func mapToContractOutputs(tx *engine.TransactionConfig) []*response.TransactionOutput {
	if tx.Outputs == nil {
		return nil
	}

	outputs := make([]*response.TransactionOutput, 0)
	for _, output := range tx.Outputs {
		outputs = append(outputs, MapToTransactionOutputContract(output))
	}
	return outputs
}

func mapToContractInputs(tx *engine.TransactionConfig) []*response.TransactionInput {
	if tx.Inputs == nil {
		return nil
	}

	inputs := make([]*response.TransactionInput, 0)
	for _, input := range tx.Inputs {
		inputs = append(inputs, MapToTransactionInputContract(input))
	}
	return inputs
}

func mapToContractIncludeUtxos(tx *engine.TransactionConfig) []*response.UtxoPointer {
	if tx.IncludeUtxos == nil {
		return nil
	}

	includeUtxos := make([]*response.UtxoPointer, 0)
	for _, utxo := range tx.IncludeUtxos {
		includeUtxos = append(includeUtxos, MapToUtxoPointer(utxo))
	}
	return includeUtxos
}

func mapToContractFromUtxos(tx *engine.TransactionConfig) []*response.UtxoPointer {
	if tx.FromUtxos == nil {
		return nil
	}

	fromUtxos := make([]*response.UtxoPointer, 0)
	for _, utxo := range tx.FromUtxos {
		fromUtxos = append(fromUtxos, MapToUtxoPointer(utxo))
	}
	return fromUtxos
}

func mapToContractDestinations(tx *engine.TransactionConfig) []*response.Destination {
	if tx.ChangeDestinations == nil {
		return nil
	}

	destinations := make([]*response.Destination, 0)
	for _, destination := range tx.ChangeDestinations {
		destinations = append(destinations, MapToDestinationContract(destination))
	}
	return destinations
}

// MapToDraftTransactionContract will map the transaction-output model from spv-wallet to the spv-wallet-models contract
func MapToDraftTransactionContract(tx *engine.DraftTransaction) *response.DraftTransaction {
	if tx == nil {
		return nil
	}

	return &response.DraftTransaction{
		Model:         *common.MapToContract(&tx.Model),
		ID:            tx.ID,
		Hex:           tx.Hex,
		XpubID:        tx.XpubID,
		ExpiresAt:     tx.ExpiresAt,
		Configuration: *MapToTransactionConfigContract(&tx.Configuration),
	}
}

// MapToTransactionInputContract will map the transaction-output model from spv-wallet-models to the spv-wallet contract
func MapToTransactionInputContract(inp *engine.TransactionInput) *response.TransactionInput {
	if inp == nil {
		return nil
	}

	return &response.TransactionInput{
		Utxo:        *MapToUtxoContract(&inp.Utxo),
		Destination: *MapToDestinationContract(&inp.Destination),
	}
}

// MapTransactionInputModelToEngine will map the transaction-output model from spv-wallet to the spv-wallet-models contract
func MapTransactionInputModelToEngine(inp *response.TransactionInput) *engine.TransactionInput {
	if inp == nil {
		return nil
	}

	return &engine.TransactionInput{
		Utxo:        *MapUtxoModelToEngine(&inp.Utxo),
		Destination: *MapDestinationModelToEngine(&inp.Destination),
	}
}

// MapToTransactionOutputContract will map the transaction-output model from spv-wallet to the spv-wallet-models contract
func MapToTransactionOutputContract(out *engine.TransactionOutput) *response.TransactionOutput {
	if out == nil {
		return nil
	}

	scriptOutputs := make([]*response.ScriptOutput, 0)
	for _, scriptOutput := range out.Scripts {
		scriptOutputs = append(scriptOutputs, MapToScriptOutputContract(scriptOutput))
	}

	return &response.TransactionOutput{
		OpReturn:     MapToOpReturnContract(out.OpReturn),
		PaymailP4:    MapToPaymailP4Contract(out.PaymailP4),
		Satoshis:     out.Satoshis,
		Script:       out.Script,
		Scripts:      scriptOutputs,
		To:           out.To,
		UseForChange: out.UseForChange,
	}
}

// MapTransactionOutputModelToEngine will map the transaction-output model from spv-wallet-models to the spv-wallet contract
func MapTransactionOutputModelToEngine(out *response.TransactionOutput) *engine.TransactionOutput {
	if out == nil {
		return nil
	}

	scriptOutputs := make([]*engine.ScriptOutput, 0)
	for _, scriptOutput := range out.Scripts {
		scriptOutputs = append(scriptOutputs, MapScriptOutputModelToEngine(scriptOutput))
	}

	return &engine.TransactionOutput{
		OpReturn:     MapOpReturnModelToEngine(out.OpReturn),
		PaymailP4:    MapPaymailP4ModelToEngine(out.PaymailP4),
		Satoshis:     out.Satoshis,
		Script:       out.Script,
		Scripts:      scriptOutputs,
		To:           out.To,
		UseForChange: out.UseForChange,
	}
}

// MapToMapProtocolContract will map the transaction-output model from spv-wallet to the spv-wallet-models contract
func MapToMapProtocolContract(mp *engine.MapProtocol) *response.MapProtocol {
	if mp == nil {
		return nil
	}

	return &response.MapProtocol{
		App:  mp.App,
		Keys: mp.Keys,
		Type: mp.Type,
	}
}

// MapMapProtocolModelToEngine will map the transaction-output model from spv-wallet-models to the spv-wallet contract
func MapMapProtocolModelToEngine(mp *response.MapProtocol) *engine.MapProtocol {
	if mp == nil {
		return nil
	}

	return &engine.MapProtocol{
		App:  mp.App,
		Keys: mp.Keys,
		Type: mp.Type,
	}
}

// MapToOpReturnContract will map the transaction-output model from spv-wallet to the spv-wallet-models contract
func MapToOpReturnContract(op *engine.OpReturn) *response.OpReturn {
	if op == nil {
		return nil
	}

	return &response.OpReturn{
		Hex:         op.Hex,
		HexParts:    op.HexParts,
		Map:         MapToMapProtocolContract(op.Map),
		StringParts: op.StringParts,
	}
}

// MapOpReturnModelToEngine will map the op-return model from spv-wallet-models to the spv-wallet contract
func MapOpReturnModelToEngine(op *response.OpReturn) *engine.OpReturn {
	if op == nil {
		return nil
	}

	return &engine.OpReturn{
		Hex:         op.Hex,
		HexParts:    op.HexParts,
		Map:         MapMapProtocolModelToEngine(op.Map),
		StringParts: op.StringParts,
	}
}
