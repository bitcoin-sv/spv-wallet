package txgraph

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/stretchr/testify/require"
)

type MinedxTxBuilder struct {
	T                    *testing.T
	P2PKHLockingScript   *script.Script
	P2PKHUnlockingScript *p2pkh.P2PKH
	HexGen               HexGen
	Block                uint32
	Satoshis             uint64
}

func (m *MinedxTxBuilder) MakeTx(inputs int) *sdk.Transaction {
	m.T.Helper()

	tx := sdk.NewTransaction()
	for i := 0; i < inputs; i++ {
		tx.AddOutput(&sdk.TransactionOutput{LockingScript: m.P2PKHLockingScript})

		vout := uint32(i)
		utxo, err := sdk.NewUTXO(m.HexGen.Val(), vout, m.P2PKHLockingScript.String(), m.Satoshis)
		require.NoError(m.T, err, "failed to initialize new utxo")
		utxo.UnlockingScriptTemplate = m.P2PKHUnlockingScript

		require.NoError(m.T, tx.AddInputsFromUTXOs(utxo), "failed to add inputs from utxo")
		m.HexGen.Inc()
	}

	merklePath := sdk.NewMerklePath(m.Block, [][]*sdk.PathElement{{{Hash: tx.TxID(), Offset: 0}}})
	m.Block++

	require.NoError(m.T, tx.AddMerkleProof(merklePath), "failed to add Merkle Proof")
	require.NoError(m.T, tx.Sign(), "failed to sign transaction")
	return tx
}

type HexGen int32

func (h *HexGen) Inc() { *h++ }

func (h HexGen) Val() string {
	return fmt.Sprintf("a%010de1b81dd2c9c0c6cd67f9bdf832e9c2bb12a1d57f30cb6ebbe78d9", h)
}
