package repository

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"iter"
	"slices"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/addresses/addressesmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"gorm.io/datatypes"
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

// Create adds a new address to the database.
func (r *Addresses) Create(ctx context.Context, newAddress *addressesmodels.NewAddress) error {
	row := &database.Address{
		UserID:             newAddress.UserID,
		Address:            newAddress.Address,
		CustomInstructions: datatypes.NewJSONSlice(newAddress.CustomInstructions),
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return spverrors.Wrapf(err, "failed to create address")
	}

	return nil
}

// FindByStringAddresses returns address rows from the database based on the provided iterator of string addresses.
func (r *Addresses) FindByStringAddresses(ctx context.Context, addresses iter.Seq[string]) ([]addressesmodels.Address, error) {
	var rows []*database.Address
	if err := r.db.
		WithContext(ctx).
		Model(&database.Address{}).
		Where("address IN ?", slices.Collect(addresses)).
		Find(&rows).Error; err != nil {
		return nil, spverrors.Wrapf(err, "failed to get addresses")
	}

	return utils.MapSlice(rows, func(row *database.Address) addressesmodels.Address {
		return addressesmodels.Address{
			Address:            row.Address,
			CreatedAt:          row.CreatedAt,
			UpdatedAt:          row.UpdatedAt,
			UserID:             row.UserID,
			CustomInstructions: (bsv.CustomInstructions)(row.CustomInstructions),
		}
	}), nil
}
