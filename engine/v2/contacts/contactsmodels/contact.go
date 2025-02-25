package contactsmodels

import (
	"time"
)

const (
	// ContactNotConfirmed is a status telling that the contact model as not confirmed yet.
	ContactNotConfirmed = "unconfirmed"
	// ContactAwaitAccept is a status telling that the contact model as invitation to add to contacts.
	ContactAwaitAccept = "awaiting"
	// ContactConfirmed is a status telling that the contact model as confirmed.
	ContactConfirmed = "confirmed"
	// ContactRejected is a status telling that the contact invitation was rejected by user.
	ContactRejected = "rejected"
)

// NewContact is a data for creating a new contact.
type NewContact struct {
	FullName          string
	RequesterPaymail  string
	NewContactPaymail string
	NewContactPubKey  string
	Status            string
	UserID            string
}

// Contact represents domain model for a user contact.
type Contact struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	FullName string
	Status   string
	Paymail  string
	PubKey   string

	UserID string
}
