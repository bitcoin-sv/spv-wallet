package repository

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"gorm.io/gorm"
	"iter"
	"slices"
)

// Addresses is a repository for addresses.
type Addresses struct {
	db *gorm.DB
}

// NewAddressesRepo creates a new repository for addresses.
func NewAddressesRepo(db *gorm.DB) *Addresses {
	return &Addresses{db: db}
}

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
