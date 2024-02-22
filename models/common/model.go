// Package common is a package that contains common models used by all other packages.
package common

import "time"

// Model is a common model that contains common fields for all models.
type Model struct {
	// CreatedAt is a time when outer model was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is a time when outer model was updated.
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt is a time when outer model was deleted.
	DeletedAt time.Time `json:"deleted_at"`
	// Metadata is a metadata map of outer model.
	Metadata map[string]interface{} `json:"metadata"`
}
