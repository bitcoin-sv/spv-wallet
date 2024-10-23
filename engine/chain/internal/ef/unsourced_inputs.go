package ef

import (
	"iter"
	"maps"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

type unsourcedInputs struct {
	// [prev txid] -> array of inputs (of current tx) that reference this prev txid
	// note that several inputs can reference the same prev txid
	prevTxToInputs map[string][]*sdk.TransactionInput
}

func findUnsourcedInputs(tx *sdk.Transaction) (*unsourcedInputs, error) {
	manager := &unsourcedInputs{
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

func (m *unsourcedInputs) addInput(input *sdk.TransactionInput) error {
	if input.SourceTXID == nil {
		return ErrMissingSourceTXID
	}
	sourceTXID := input.SourceTXID.String()
	m.prevTxToInputs[sourceTXID] = append(m.prevTxToInputs[sourceTXID], input)
	return nil
}

func (m *unsourcedInputs) txCount() int {
	return len(m.prevTxToInputs)
}

func (m *unsourcedInputs) getMissingTXIDs() iter.Seq[string] {
	return maps.Keys(m.prevTxToInputs)
}

func (m *unsourcedInputs) hydrate(sourceTx *sdk.Transaction) error {
	sourceTXID := sourceTx.TxID().String()
	inputs, ok := m.prevTxToInputs[sourceTXID]
	if !ok {
		return spverrors.Newf("got not requested transaction: %s", sourceTXID)
	}
	for _, input := range inputs {
		input.SourceTransaction = sourceTx
	}
	delete(m.prevTxToInputs, sourceTXID)
	return nil
}
