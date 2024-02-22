// Package models contains all models (contracts) between spv-wallet api and other spv-wallet solutions
package models

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/common"
)

// AccessKey is a model that represents an access key.
type AccessKey struct {
	// Model is a common model that contains common fields for all models.
	common.Model

	// ID is an access key id.
	ID string `json:"id"`
	// XpubID is an access key's xpub related id.
	XpubID string `json:"xpub_id"`
	// RevokedAt is a time when access key was revoked.
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	// Key is a string representation of an access key.
	Key string `json:"key,omitempty"`
}
