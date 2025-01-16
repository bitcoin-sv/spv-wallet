package repository

import (
	"context"
	"iter"
	"slices"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"gorm.io/gorm"
)

// Addresses is a repository for addresses.
type Addresses struct {
	db *gorm.DB
}

// NewAddressesRepo creates a new repository for addresses.
func NewAddressesRepo(db *gorm.DB) *Addresses {
	return &Addresses{db: db}
}

// FindByStringAddresses returns address rows from the database based on the provided iterator of string addresses.
func (r *Addresses) FindByStringAddresses(ctx context.Context, addresses iter.Seq[string]) ([]*database.Address, error) {
	var rows []*database.Address
	if err := r.db.
		WithContext(ctx).
		Model(&database.Address{}).
		Where("address IN ?", slices.Collect(addresses)).
		Find(&rows).Error; err != nil {
		return nil, spverrors.Wrapf(err, "failed to get addresses")
	}

	return rows, nil
}
