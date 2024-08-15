package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapToOldTransactionContract will map the model from spv-wallet to the spv-wallet-models contract
func MapToOldTransactionContract(t *engine.Transaction) *models.Transaction {
	if t == nil {
		return nil
	}

	model := models.Transaction{
		Model:                *common.MapToOldContract(&t.Model),
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

	processOldMetadata(t, t.XPubID, &model)
	processOldOutputValue(t, t.XPubID, &model)

	return &model
}

// MapToOldTransactionContractForAdmin will map the model from spv-wallet to the spv-wallet-models contract for admin
func MapToOldTransactionContractForAdmin(t *engine.Transaction) *models.Transaction {
	if t == nil {
		return nil
	}

	model := models.Transaction{
		Model:           *common.MapToOldContract(&t.Model),
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

	processOldMetadata(t, t.XPubID, &model)

	return &model
}

func processOldMetadata(t *engine.Transaction, xpubID string, model *models.Transaction) {
	if len(t.XpubMetadata) > 0 && len(t.XpubMetadata[xpubID]) > 0 {
		if t.Model.Metadata == nil {
			model.Model.Metadata = make(models.Metadata)
		}
		for key, value := range t.XpubMetadata[xpubID] {
			model.Model.Metadata[key] = value
		}
	}
}

func processOldOutputValue(t *engine.Transaction, xpubID string, model *models.Transaction) {
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

// MapOldTransactionModelToEngine will map the model from spv-wallet-models to the spv-wallet contract
func MapOldTransactionModelToEngine(t *models.Transaction) *engine.Transaction {
	if t == nil {
		return nil
	}

	return &engine.Transaction{
		Model:           *common.MapOldContractToModel(&t.Model),
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

// MapOldTransactionConfigEngineToModel will map the transaction-config model from spv-wallet to the spv-wallet-models contract
func MapOldTransactionConfigEngineToModel(tx *models.TransactionConfig) *engine.TransactionConfig {
	if tx == nil {
		return nil
	}

	return &engine.TransactionConfig{
		ChangeDestinations:         mapToOldEngineDestinations(tx),
		ChangeDestinationsStrategy: engine.ChangeStrategy(tx.ChangeStrategy),
		ChangeMinimumSatoshis:      tx.ChangeMinimumSatoshis,
		ChangeNumberOfDestinations: tx.ChangeNumberOfDestinations,
		ChangeSatoshis:             tx.ChangeSatoshis,
		ExpiresIn:                  tx.ExpiresIn,
		Fee:                        tx.Fee,
		FeeUnit:                    MapOldFeeUnitModelToEngine(tx.FeeUnit),
		FromUtxos:                  mapToOldEngineFromUtxos(tx),
		IncludeUtxos:               mapOldIncludeUtxosModelToEngine(tx),
		Inputs:                     mapToOldEngineInputs(tx),
		Outputs:                    mapToOldEngineOutputs(tx),
		SendAllTo:                  MapOldTransactionOutputModelToEngine(tx.SendAllTo),
		Sync:                       MapOldSyncConfigModelToEngine(tx.Sync),
	}
}

func mapToOldEngineOutputs(tx *models.TransactionConfig) []*engine.TransactionOutput {
	if tx.Outputs == nil {
		return nil
	}

	outputs := make([]*engine.TransactionOutput, 0)
	for _, output := range tx.Outputs {
		outputs = append(outputs, MapOldTransactionOutputModelToEngine(output))
	}
	return outputs
}

func mapToOldEngineInputs(tx *models.TransactionConfig) []*engine.TransactionInput {
	if tx.Inputs == nil {
		return nil
	}

	inputs := make([]*engine.TransactionInput, 0)
	for _, input := range tx.Inputs {
		inputs = append(inputs, MapOldTransactionInputModelToEngine(input))
	}
	return inputs
}

func mapOldIncludeUtxosModelToEngine(tx *models.TransactionConfig) []*engine.UtxoPointer {
	if tx.IncludeUtxos == nil {
		return nil
	}

	includeUtxos := make([]*engine.UtxoPointer, 0)
	for _, utxo := range tx.IncludeUtxos {
		includeUtxos = append(includeUtxos, MapOldUtxoPointerModelToEngine(utxo))
	}
	return includeUtxos
}

func mapToOldEngineFromUtxos(tx *models.TransactionConfig) []*engine.UtxoPointer {
	if tx.FromUtxos == nil {
		return nil
	}

	fromUtxos := make([]*engine.UtxoPointer, 0)
	for _, utxo := range tx.FromUtxos {
		fromUtxos = append(fromUtxos, MapOldUtxoPointerModelToEngine(utxo))
	}
	return fromUtxos
}

func mapToOldEngineDestinations(tx *models.TransactionConfig) []*engine.Destination {
	if tx.ChangeDestinations == nil {
		return nil
	}

	destinations := make([]*engine.Destination, 0)
	for _, destination := range tx.ChangeDestinations {
		destinations = append(destinations, MapOldDestinationModelToEngine(destination))
	}
	return destinations
}

// MapToOldTransactionConfigContract will map the transaction-config model from spv-wallet-models to the spv-wallet contract
func MapToOldTransactionConfigContract(tx *engine.TransactionConfig) *models.TransactionConfig {
	if tx == nil {
		return nil
	}

	return &models.TransactionConfig{
		ChangeDestinations:         mapToOldContractDestinations(tx),
		ChangeStrategy:             string(tx.ChangeDestinationsStrategy),
		ChangeMinimumSatoshis:      tx.ChangeMinimumSatoshis,
		ChangeNumberOfDestinations: tx.ChangeNumberOfDestinations,
		ChangeSatoshis:             tx.ChangeSatoshis,
		ExpiresIn:                  tx.ExpiresIn,
		FeeUnit:                    MapToOldFeeUnitContract(tx.FeeUnit),
		FromUtxos:                  mapToOldContractFromUtxos(tx),
		IncludeUtxos:               mapToOldContractIncludeUtxos(tx),
		Inputs:                     mapToOldContractInputs(tx),
		Outputs:                    mapToOldContractOutputs(tx),
		SendAllTo:                  MapToOldTransactionOutputContract(tx.SendAllTo),
		Sync:                       MapToOldSyncConfigContract(tx.Sync),
	}
}

func mapToOldContractOutputs(tx *engine.TransactionConfig) []*models.TransactionOutput {
	if tx.Outputs == nil {
		return nil
	}

	outputs := make([]*models.TransactionOutput, 0)
	for _, output := range tx.Outputs {
		outputs = append(outputs, MapToOldTransactionOutputContract(output))
	}
	return outputs
}

func mapToOldContractInputs(tx *engine.TransactionConfig) []*models.TransactionInput {
	if tx.Inputs == nil {
		return nil
	}

	inputs := make([]*models.TransactionInput, 0)
	for _, input := range tx.Inputs {
		inputs = append(inputs, MapToOldTransactionInputContract(input))
	}
	return inputs
}

func mapToOldContractIncludeUtxos(tx *engine.TransactionConfig) []*models.UtxoPointer {
	if tx.IncludeUtxos == nil {
		return nil
	}

	includeUtxos := make([]*models.UtxoPointer, 0)
	for _, utxo := range tx.IncludeUtxos {
		includeUtxos = append(includeUtxos, MapToOldUtxoPointer(utxo))
	}
	return includeUtxos
}

func mapToOldContractFromUtxos(tx *engine.TransactionConfig) []*models.UtxoPointer {
	if tx.FromUtxos == nil {
		return nil
	}

	fromUtxos := make([]*models.UtxoPointer, 0)
	for _, utxo := range tx.FromUtxos {
		fromUtxos = append(fromUtxos, MapToOldUtxoPointer(utxo))
	}
	return fromUtxos
}

func mapToOldContractDestinations(tx *engine.TransactionConfig) []*models.Destination {
	if tx.ChangeDestinations == nil {
		return nil
	}

	destinations := make([]*models.Destination, 0)
	for _, destination := range tx.ChangeDestinations {
		destinations = append(destinations, MapOldToDestinationContract(destination))
	}
	return destinations
}

// MapToOldDraftTransactionContract will map the transaction-output model from spv-wallet to the spv-wallet-models contract
func MapToOldDraftTransactionContract(tx *engine.DraftTransaction) *models.DraftTransaction {
	if tx == nil {
		return nil
	}

	return &models.DraftTransaction{
		Model:         *common.MapToOldContract(&tx.Model),
		ID:            tx.ID,
		Hex:           tx.Hex,
		XpubID:        tx.XpubID,
		ExpiresAt:     tx.ExpiresAt,
		Configuration: *MapToOldTransactionConfigContract(&tx.Configuration),
	}
}

// MapToOldTransactionInputContract will map the transaction-output model from spv-wallet-models to the spv-wallet contract
func MapToOldTransactionInputContract(inp *engine.TransactionInput) *models.TransactionInput {
	if inp == nil {
		return nil
	}

	return &models.TransactionInput{
		Utxo:        *MapToOldUtxoContract(&inp.Utxo),
		Destination: *MapOldToDestinationContract(&inp.Destination),
	}
}

// MapOldTransactionInputModelToEngine will map the transaction-output model from spv-wallet to the spv-wallet-models contract
func MapOldTransactionInputModelToEngine(inp *models.TransactionInput) *engine.TransactionInput {
	if inp == nil {
		return nil
	}

	return &engine.TransactionInput{
		Utxo:        *MapOldUtxoModelToEngine(&inp.Utxo),
		Destination: *MapOldDestinationModelToEngine(&inp.Destination),
	}
}

// MapToOldTransactionOutputContract will map the transaction-output model from spv-wallet to the spv-wallet-models contract
func MapToOldTransactionOutputContract(out *engine.TransactionOutput) *models.TransactionOutput {
	if out == nil {
		return nil
	}

	scriptOutputs := make([]*models.ScriptOutput, 0)
	for _, scriptOutput := range out.Scripts {
		scriptOutputs = append(scriptOutputs, MapToOldScriptOutputContract(scriptOutput))
	}

	return &models.TransactionOutput{
		OpReturn:     MapToOldOpReturnContract(out.OpReturn),
		PaymailP4:    MapToOldPaymailP4Contract(out.PaymailP4),
		Satoshis:     out.Satoshis,
		Script:       out.Script,
		Scripts:      scriptOutputs,
		To:           out.To,
		UseForChange: out.UseForChange,
	}
}

// MapOldTransactionOutputModelToEngine will map the transaction-output model from spv-wallet-models to the spv-wallet contract
func MapOldTransactionOutputModelToEngine(out *models.TransactionOutput) *engine.TransactionOutput {
	if out == nil {
		return nil
	}

	scriptOutputs := make([]*engine.ScriptOutput, 0)
	for _, scriptOutput := range out.Scripts {
		scriptOutputs = append(scriptOutputs, MapOldScriptOutputModelToEngine(scriptOutput))
	}

	return &engine.TransactionOutput{
		OpReturn:     MapOldOpReturnModelToEngine(out.OpReturn),
		PaymailP4:    MapOldPaymailP4ModelToEngine(out.PaymailP4),
		Satoshis:     out.Satoshis,
		Script:       out.Script,
		Scripts:      scriptOutputs,
		To:           out.To,
		UseForChange: out.UseForChange,
	}
}

// MapToOldMapProtocolContract will map the transaction-output model from spv-wallet to the spv-wallet-models contract
func MapToOldMapProtocolContract(mp *engine.MapProtocol) *models.MapProtocol {
	if mp == nil {
		return nil
	}

	return &models.MapProtocol{
		App:  mp.App,
		Keys: mp.Keys,
		Type: mp.Type,
	}
}

// MapOldMapProtocolModelToEngine will map the transaction-output model from spv-wallet-models to the spv-wallet contract
func MapOldMapProtocolModelToEngine(mp *models.MapProtocol) *engine.MapProtocol {
	if mp == nil {
		return nil
	}

	return &engine.MapProtocol{
		App:  mp.App,
		Keys: mp.Keys,
		Type: mp.Type,
	}
}

// MapToOldOpReturnContract will map the transaction-output model from spv-wallet to the spv-wallet-models contract
func MapToOldOpReturnContract(op *engine.OpReturn) *models.OpReturn {
	if op == nil {
		return nil
	}

	return &models.OpReturn{
		Hex:         op.Hex,
		HexParts:    op.HexParts,
		Map:         MapToOldMapProtocolContract(op.Map),
		StringParts: op.StringParts,
	}
}

// MapOldOpReturnModelToEngine will map the op-return model from spv-wallet-models to the spv-wallet contract
func MapOldOpReturnModelToEngine(op *models.OpReturn) *engine.OpReturn {
	if op == nil {
		return nil
	}

	return &engine.OpReturn{
		Hex:         op.Hex,
		HexParts:    op.HexParts,
		Map:         MapOldMapProtocolModelToEngine(op.Map),
		StringParts: op.StringParts,
	}
}
