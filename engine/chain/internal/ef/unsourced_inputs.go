package ef

import (
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"iter"
	"maps"
)

type unsourcedInputsManager struct {
	// [prev txid] -> array of inputs (of current tx) that reference this prev txid
	// note that several inputs can reference the same prev txid
	prevTxToInputs map[string][]*sdk.TransactionInput
}

func findUnsourcedInputs(tx *sdk.Transaction) (*unsourcedInputsManager, error) {
	manager := &unsourcedInputsManager{
		prevTxToInputs: make(map[string][]*sdk.TransactionInput),
	}

	for _, input := range tx.Inputs {
		if input.SourceTransaction != nil {
			continue
		}
		if err := manager.addInput(input); err != nil {
			return manager, err
		}
	}

	return manager, nil
}

func (m *unsourcedInputsManager) addInput(input *sdk.TransactionInput) error {
	if input.SourceTXID == nil {
		return ErrMissingSourceTXID
	}
	sourceTXID := input.SourceTXID.String()
	m.prevTxToInputs[sourceTXID] = append(m.prevTxToInputs[sourceTXID], input)
	return nil
}

func (m *unsourcedInputsManager) txCount() int {
	return len(m.prevTxToInputs)
}

func (m *unsourcedInputsManager) getMissingTXIDs() iter.Seq[string] {
	return maps.Keys(m.prevTxToInputs)
}

func (m *unsourcedInputsManager) deleteTXID(txid string) {
	delete(m.prevTxToInputs, txid)
}

func (m *unsourcedInputsManager) hydrate(sourceTx *sdk.Transaction) error {
	sourceTXID := sourceTx.TxID().String()
	inputs, ok := m.prevTxToInputs[sourceTXID]
	if !ok {
		return spverrors.Newf("got not requested transaction: %s", sourceTXID)
	}
	for _, input := range inputs {
		input.SourceTransaction = sourceTx
	}
	return nil
}
