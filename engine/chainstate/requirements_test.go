package chainstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkRequirement(t *testing.T) {
	t.Parallel()

	t.Run("found in mempool - mAPI", func(t *testing.T) {
		success := checkRequirementMapi(requiredInMempool, onChainExample1TxID, &TransactionInfo{
			BlockHash:     "",
			BlockHeight:   0,
			Confirmations: 0,
			ID:            onChainExample1TxID,
			MinerID:       minerTaal.MinerID,
			Provider:      "",
		})
		assert.Equal(t, true, success)
	})

	t.Run("found in mempool - on-chain - mAPI", func(t *testing.T) {
		success := checkRequirementMapi(requiredInMempool, onChainExample1TxID, &TransactionInfo{
			BlockHash:     onChainExample1BlockHash,
			BlockHeight:   onChainExample1BlockHeight,
			Confirmations: 1,
			ID:            onChainExample1TxID,
			MinerID:       minerTaal.MinerID,
			Provider:      "",
		})
		assert.Equal(t, true, success)
	})

	t.Run("found in mempool - whatsonchain", func(t *testing.T) {
		success := checkRequirementMapi(requiredInMempool, onChainExample1TxID, &TransactionInfo{
			BlockHash:     "",
			BlockHeight:   0,
			Confirmations: 0,
			ID:            onChainExample1TxID,
			MinerID:       "",
			Provider:      "whatsonchain",
		})
		assert.Equal(t, true, success)
	})

	t.Run("not in mempool - mAPI", func(t *testing.T) {
		success := checkRequirementMapi(requiredInMempool, onChainExample1TxID, &TransactionInfo{
			BlockHash:     "",
			BlockHeight:   0,
			Confirmations: 0,
			ID:            "",
			MinerID:       minerTaal.MinerID,
			Provider:      "",
		})
		assert.Equal(t, false, success)
	})

	t.Run("found on chain - mAPI", func(t *testing.T) {
		success := checkRequirementMapi(requiredOnChain, onChainExample1TxID, &TransactionInfo{
			BlockHash:     onChainExample1BlockHash,
			BlockHeight:   onChainExample1BlockHeight,
			Confirmations: 1,
			ID:            onChainExample1TxID,
			MinerID:       minerTaal.MinerID,
			Provider:      "",
		})
		assert.Equal(t, true, success)
	})

	t.Run("not on chain - mAPI", func(t *testing.T) {
		success := checkRequirementMapi(requiredOnChain, onChainExample1TxID, &TransactionInfo{
			BlockHash:     "",
			BlockHeight:   0,
			Confirmations: 0,
			ID:            onChainExample1TxID,
			MinerID:       minerTaal.MinerID,
			Provider:      "",
		})
		assert.Equal(t, false, success)
	})
}
