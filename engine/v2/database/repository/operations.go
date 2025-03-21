package repository

import (
	"context"
	"iter"
	"slices"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database/dbquery"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/operations/operationsmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Operations is a repository for operations.
type Operations struct {
	db *gorm.DB
}

// NewOperationsRepo creates a new repository for operations.
func NewOperationsRepo(db *gorm.DB) *Operations {
	return &Operations{db: db}
}

// PaginatedForUser returns operations for a user based on userID and the provided paging options.
func (o *Operations) PaginatedForUser(ctx context.Context, userID string, page filter.Page) (*models.PagedResult[operationsmodels.Operation], error) {
	rows, err := dbquery.PaginatedQuery[database.Operation](
		ctx,
		page,
		o.db,
		dbquery.UserID(userID),
		dbquery.Preload("Transaction"),
	)
	if err != nil {
		return nil, err
	}
	return &models.PagedResult[operationsmodels.Operation]{
		PageDescription: rows.PageDescription,
		Content: lo.Map(rows.Content, func(operation *database.Operation, _ int) *operationsmodels.Operation {
			return &operationsmodels.Operation{
				TxID:         operation.TxID,
				UserID:       operation.UserID,
				CreatedAt:    operation.CreatedAt,
				Counterparty: operation.Counterparty,
				Type:         operation.Type,
				Value:        operation.Value,
				TxStatus:     operation.Transaction.TxStatus,
				BlockHeight:  operation.Transaction.BlockHeight,
				BlockHash:    operation.Transaction.BlockHash,
			}
		}),
	}, nil
}

// SaveAll saves operations to the database.
func (o *Operations) SaveAll(ctx context.Context, operations iter.Seq[*txmodels.NewOperation]) error {
	rows := mapOperations(operations)

	query := o.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		})

	if err := query.Create(rows).Error; err != nil {
		return err
	}

	return nil
}

func mapOperations(operations iter.Seq[*txmodels.NewOperation]) []database.Operation {
	transactions := map[string]*database.TrackedTransaction{}
	var tx *database.TrackedTransaction
	var ok bool

	return slices.Collect(func(yield func(database.Operation) bool) {
		for operation := range operations {
			tx, ok = transactions[operation.Transaction.ID]
			if !ok {
				tx = mapTransaction(operation)
				transactions[operation.Transaction.ID] = tx
			}

			yield(database.Operation{
				UserID: operation.UserID,

				Counterparty: operation.Counterparty,
				Type:         operation.Type,
				Value:        operation.Value,

				TxID:        operation.Transaction.ID,
				Transaction: tx,
			})
		}
	})
}

func mapTransaction(operation *txmodels.NewOperation) *database.TrackedTransaction {
	beefHex := operation.Transaction.BEEFHex()
	rawHex := operation.Transaction.RawHex()
	tx := &database.TrackedTransaction{
		ID:       operation.Transaction.ID,
		TxStatus: string(operation.Transaction.TxStatus),
		BeefHex:  lo.If(beefHex != "", &beefHex).Else(nil),
		RawHex:   lo.If(rawHex != "", &rawHex).Else(nil),
	}

	for _, input := range operation.Transaction.TransactionInputSources() {
		tx.SourceTxInputs = append(tx.SourceTxInputs, database.TxInput{
			TxID:       input.TxID,
			SourceTxID: input.SourceTxID,
		})
	}

	for _, input := range operation.Transaction.Inputs {
		tx.Inputs = append(tx.Inputs, &database.TrackedOutput{
			TxID:       input.TxID,
			Vout:       input.Vout,
			SpendingTX: operation.Transaction.ID,
			UserID:     input.UserID,
			Satoshis:   input.Satoshis,
		})
	}

	for _, output := range operation.Transaction.Outputs {
		if output.UTXO != nil {
			tx.CreateUTXO(
				&database.TrackedOutput{
					TxID:     operation.Transaction.ID,
					Vout:     output.Vout,
					UserID:   output.UserID,
					Satoshis: output.Satoshis,
				},
				output.Bucket,
				output.UTXO.EstimatedInputSize,
				output.UTXO.CustomInstructions,
			)
		} else if output.Data != nil {
			tx.CreateDataOutput(&database.Data{
				TxID:   operation.Transaction.ID,
				Vout:   output.Vout,
				UserID: output.UserID,
				Blob:   output.Data,
			})
		}
	}

	return tx
}
