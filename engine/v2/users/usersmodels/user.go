package usersmodels

import (
	"fmt"
	"time"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
)

// NewUser represents data for creating a new user
type NewUser struct {
	PublicKey string
	Paymail   *NewPaymail
}

// User is a domain model for existing users
type User struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time

	PublicKey string
	Paymails  []*Paymail
}

// PubKeyObj returns the go-sdk primitives.PublicKey object from the user's PubKey string
func (u *User) PubKeyObj() (*primitives.PublicKey, error) {
	pub, err := primitives.PublicKeyFromString(u.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("invalid PubKey: %w", err)
	}
	return pub, nil
}
