package database

import (
	"slices"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/beef"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/gorm"
)

// TrackedTransaction represents a transaction in the database.
type TrackedTransaction struct {
	ID       string `gorm:"type:char(64);primaryKey"`
	TxStatus string

	CreatedAt time.Time
	UpdatedAt time.Time

	Data []*Data `gorm:"foreignKey:TxID"`

	Inputs  []*TrackedOutput `gorm:"foreignKey:SpendingTX"`
	Outputs []*TrackedOutput `gorm:"foreignKey:TxID"`

	newUTXOs []*UserUTXO `gorm:"-"`

	BeefHex            string               `gorm:"column:beef_hex"`
	RawHex             string               `gorm:"column:raw_hex"`
	SourceTransactions []TrackedTransaction `gorm:"many2many:source_transactions"`
}

// ToTxQueryResult converts a TrackedTransaction into a TxQueryResult.
func (t *TrackedTransaction) ToTxQueryResult() *beef.TxQueryResult {
	return &beef.TxQueryResult{
		SourceTXID: t.ID,
		RawHex:     &t.RawHex,
		BeefHex:    &t.BeefHex,
	}
}

// HasBeefHex checks if the tracked transaction record contains a non-empty BeefHex attribute.
func (t *TrackedTransaction) HasBeefHex() bool { return t.BeefHex != "" }

// HasRawHex checks if the tracked transaction record does not contain a BeefHex attribute.
func (t *TrackedTransaction) HasRawHex() bool { return !t.HasBeefHex() }

// CreateUTXO prepares a new UTXO and adds it to the transaction.
func (t *TrackedTransaction) CreateUTXO(
	output *TrackedOutput,
	bucket string,
	estimatedInputSize uint64,
	customInstructions bsv.CustomInstructions,
) {
	t.Outputs = append(t.Outputs, output)
	t.newUTXOs = append(t.newUTXOs, NewUTXO(output, bucket, estimatedInputSize, customInstructions))
}

// CreateDataOutput prepares a new Data output and adds it to the transaction.
func (t *TrackedTransaction) CreateDataOutput(data *Data) {
	t.Data = append(t.Data, data)
}

// AfterCreate is a hook that is called after creating the transaction.
// It is responsible for adding new (User's) UTXOs and removing spent UTXOs.
func (t *TrackedTransaction) AfterCreate(tx *gorm.DB) error {
	// Add new UTXOs
	if len(t.newUTXOs) > 0 {
		err := tx.Model(&UserUTXO{}).Create(t.newUTXOs).Error
		if err != nil {
			return spverrors.Wrapf(err, "failed to save user utxos")
		}
	}

	spentOutpoints := slices.AppendSeq(
		make([][]any, 0, len(t.Inputs)),
		func(yield func(sqlPair []any) bool) {
			for _, outpoint := range t.Inputs {
				yield([]any{outpoint.TxID, outpoint.Vout})
			}
		})
	if len(spentOutpoints) > 0 {
		// Remove spent UTXOs
		err := tx.
			Where("(tx_id, vout) IN ?", spentOutpoints).
			Delete(&UserUTXO{}).
			Error
		if err != nil {
			return spverrors.Wrapf(err, "failed to delete spent utxos")
		}
	}

	return nil
}
