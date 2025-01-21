package domainmodels

import (
	"time"
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
}

// NewPaymail represents data for creating a new paymail
type NewPaymail struct {
	Alias      string
	Domain     string
	PublicName string
	Avatar     string
	UserID     string
}
