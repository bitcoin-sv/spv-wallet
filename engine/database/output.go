package database

import "github.com/bitcoin-sv/spv-wallet/models/bsv"

type Output interface {
	IsSpent() bool
	Outpoint() *bsv.Outpoint

	ToTrackedOutput() *TrackedOutput
	ToUserUTXO() *UserUtxos
}

func NewDataOutput(txID string, vout uint32) Output {
	return &virtualOutput{
		txID:   txID,
		vout:   vout,
		bucket: "data",
	}
}

func NewP2PKHOutput(txID string, vout uint32, userID string, satoshis bsv.Satoshis) Output {
	return &virtualOutput{
		txID:                         txID,
		vout:                         vout,
		bucket:                       "p2pkh", //TODO: check if this is correct
		userID:                       userID,
		satoshis:                     satoshis,
		unlockingScriptEstimatedSize: 106, //TODO: check if this is correct
	}
}

type virtualOutput struct {
	txID       string
	vout       uint32
	spendingTX string

	bucket                       string
	satoshis                     bsv.Satoshis
	unlockingScriptEstimatedSize uint64

	userID string
}

func (o *virtualOutput) IsSpent() bool {
	return o.spendingTX != ""
}

func (o *virtualOutput) Outpoint() *bsv.Outpoint {
	return &bsv.Outpoint{
		TxID: o.txID,
		Vout: o.vout,
	}
}

func (o *virtualOutput) ToTrackedOutput() *TrackedOutput {
	return &TrackedOutput{
		TxID:       o.txID,
		Vout:       o.vout,
		SpendingTX: o.spendingTX,
	}
}

func (o *virtualOutput) ToUserUTXO() *UserUtxos {
	return &UserUtxos{
		UserID:                       o.userID,
		TxID:                         o.txID,
		Vout:                         o.vout,
		Satoshis:                     uint64(o.satoshis),
		UnlockingScriptEstimatedSize: o.unlockingScriptEstimatedSize,
		Bucket:                       o.bucket,
	}
}
