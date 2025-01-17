package usermodels

import (
	"time"

	paymailmodels "github.com/bitcoin-sv/spv-wallet/engine/paymail/models"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// User is a domain model for existing users
type User struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time

	PublicKey string
	Paymails  []*paymailmodels.Paymail
}

// NewAddress represents data for creating a new address
type NewAddress struct {
	Address            string
	CustomInstructions bsv.CustomInstructions
}

// NewUser represents data for creating a new user
type NewUser struct {
	PublicKey string
	Paymail   *NewPaymail
}

// NewPaymail represents data for creating a new paymail
type NewPaymail struct {
	Alias      string
	Domain     string
	PublicName string
	Avatar     string
}
