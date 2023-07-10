// Package common is a package that contains common models & methods used by all other packages.
package common

import (
	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-models/common"
)

// MapToContract will map the common model to the bux-models contract
func MapToContract(m *bux.Model) *common.Model {
	return &common.Model{
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt.Time,
		Metadata:  m.Metadata,
	}
}

// MapToModel will map the bux-models contract to the common bux model
func MapToModel(m *common.Model) *bux.Model {
	return &bux.Model{
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Metadata:  m.Metadata,
	}
}
