package engine

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/libsv/go-bt/v2"
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

	ver := bt.LittleEndianBytes(beefTx.version, 4)
	ver[2] = 0xBE
	ver[3] = 0xEF
	beefSize += len(ver)

	nBUMPS := bt.VarInt(len(beefTx.bumps)).Bytes()
	beefSize += len(nBUMPS)

	bumps := beefTx.bumps.Bytes()
	beefSize += len(bumps)

	nTransactions := bt.VarInt(uint64(len(beefTx.transactions))).Bytes()
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

func toBeefBytes(tx *bt.Tx, bumps BUMPs) []byte {
	txBeefBytes := tx.Bytes()

	bumpIdx := getBumpPathIndex(tx, bumps)
	if bumpIdx > -1 {
		txBeefBytes = append(txBeefBytes, hasBUMP)
		txBeefBytes = append(txBeefBytes, bt.VarInt(bumpIdx).Bytes()...)
	} else {
		txBeefBytes = append(txBeefBytes, hasNoBUMP)
	}

	return txBeefBytes
}

func getBumpPathIndex(tx *bt.Tx, bumps BUMPs) int {
	bumpIndex := -1

	for i, bump := range bumps {
		for _, path := range bump.Path[0] {
			if path.Hash == tx.TxID() {
				bumpIndex = i
			}
		}
	}

	return bumpIndex
}
