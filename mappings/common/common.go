// Package common is a package that contains common models & methods used by all other packages.
package common

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models/common"
)

// MapToOldContract will map the common model to the spv-wallet-models contract (this is deprecated and will be removed in the future)
func MapToOldContract(m *engine.Model) *common.OldModel {
	if m == nil {
		return nil
	}

	result := common.OldModel{
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Metadata:  m.Metadata,
	}
	if m.DeletedAt.Valid {
		result.DeletedAt = &m.DeletedAt.Time
	}

	return &result
}

// MapToContract will map the common model to the spv-wallet-models contract
func MapToContract(m *engine.Model) *common.Model {
	if m == nil {
		return nil
	}

	result := common.Model{
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Metadata:  m.Metadata,
	}
	if m.DeletedAt.Valid {
		result.DeletedAt = &m.DeletedAt.Time
	}

	return &result
}

// MapOldContractToModel will map the spv-wallet-models contract to the common SPV Wallet Model (this is deprecated and will be removed in the future)
func MapOldContractToModel(m *common.OldModel) *engine.Model {
	if m == nil {
		return nil
	}

	return &engine.Model{
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Metadata:  m.Metadata,
	}
}

// MapToModel will map the spv-wallet-models contract to the common SPV Wallet Model
func MapToModel(m *common.Model) *engine.Model {
	if m == nil {
		return nil
	}

	return &engine.Model{
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Metadata:  m.Metadata,
	}
}
