package engine

import (
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/util"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

var (
	hasBUMP   = byte(0x01)
	hasNoBUMP = byte(0x00)
)

func (beefTx *beefTx) toBeefBytes() ([]byte, error) {
	if len(beefTx.bumps) == 0 || len(beefTx.transactions) < 2 { // valid BEEF contains at least two transactions (new transaction and one parent transaction)
		return nil, spverrors.Newf("beef tx is incomplete")
	}

	// get beef bytes
	beefSize := 0

	ver := util.LittleEndianBytes(beefTx.version, 4)
	ver[2] = 0xBE
	ver[3] = 0xEF
	beefSize += len(ver)

	nBUMPS := trx.VarInt(len(beefTx.bumps)).Bytes()
	beefSize += len(nBUMPS)

	bumps := beefTx.bumps.Bytes()
	beefSize += len(bumps)

	nTransactions := trx.VarInt(uint64(len(beefTx.transactions))).Bytes()
	beefSize += len(nTransactions)

	transactions := make([][]byte, 0, len(beefTx.transactions))

	for _, t := range beefTx.transactions {
		txBytes := toBeefBytes(t, beefTx.bumps)
		transactions = append(transactions, txBytes)
		beefSize += len(txBytes)
	}

	// compose beef
	buffer := make([]byte, 0, beefSize)
	buffer = append(buffer, ver...)
	buffer = append(buffer, nBUMPS...)
	buffer = append(buffer, bumps...)

	buffer = append(buffer, nTransactions...)

	for _, t := range transactions {
		buffer = append(buffer, t...)
	}

	return buffer, nil
}

func toBeefBytes(tx *trx.Transaction, bumps BUMPs) []byte {
	txBeefBytes := tx.Bytes()

	bumpIdx, ok := getBumpPathIndex(tx, bumps)
	if ok {
		txBeefBytes = append(txBeefBytes, hasBUMP)
		txBeefBytes = append(txBeefBytes, trx.VarInt(bumpIdx).Bytes()...)
	} else {
		txBeefBytes = append(txBeefBytes, hasNoBUMP)
	}

	return txBeefBytes
}

func getBumpPathIndex(tx *trx.Transaction, bumps BUMPs) (uint64, bool) {
	bumpIndex := uint64(0)
	found := false
	txID := tx.TxID().String()

	for i := uint64(0); i < uint64(len(bumps)); i++ {
		for _, path := range bumps[i].Path[0] {
			if path.Hash == txID {
				bumpIndex = i
				found = true //TODO should we break here? (just return the first one or browse to return the last one???)
			}
		}
	}

	if !found {
		return 0, false
	}

	return bumpIndex, true
}
