package models

import (
	"time"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
)

// Paymail is a domain model for existing paymail
type Paymail struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	// TODO: Handle DeletedAt

	Alias  string
	Domain string

	PublicName string
	Avatar     string

	UserID string
	User   User
}

// User represents a user interface
// NOTE: Cannot used usermodels.User directly because of circular dependency
type User interface {
	PubKeyObj() (*primitives.PublicKey, error)
}
