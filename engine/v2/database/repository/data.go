package repository

import (
	"context"
	"errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/data/datamodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/gorm"
)

// Data is a repository for data.
type Data struct {
	db *gorm.DB
}

// NewDataRepo creates a new instance of the data repository.
func NewDataRepo(db *gorm.DB) *Data {
	return &Data{
		db: db,
	}
}

// FindForUser returns the data by outpoint for a specific user.
func (r *Data) FindForUser(ctx context.Context, id string, userID string) (*datamodels.Data, error) {
	outpoint, err := bsv.OutpointFromString(id)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to parse Data ID to outpoint")
	}

	var row database.Data
	if err := r.db.WithContext(ctx).
		Where("tx_id = ? AND vout = ? AND user_id = ?", outpoint.TxID, outpoint.Vout, userID).
		First(&row).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &datamodels.Data{
		TxID:   row.TxID,
		Vout:   row.Vout,
		UserID: row.UserID,
		Blob:   row.Blob,
	}, nil

}
