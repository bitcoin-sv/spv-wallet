package repository

import (
	"context"
	"errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/data/datamodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/gorm"
)

type Data struct {
	db *gorm.DB
}

func NewDataRepo(db *gorm.DB) *Data {
	return &Data{
		db: db,
	}
}

func (r *Data) FindForUser(ctx context.Context, outpoint bsv.Outpoint, userID string) (*datamodels.Data, error) {
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
