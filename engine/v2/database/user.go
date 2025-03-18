package database

import (
	"fmt"
	"time"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	"gorm.io/gorm"
)

// User represents a user in the database
type User struct {
	ID        string `gorm:"type:char(34);primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	PubKey string `gorm:"index;unique;not null"`

	Paymails  []*Paymail     `gorm:"foreignKey:UserID"`
	Addresses []*Address     `gorm:"foreignKey:UserID"`
	Contacts  []*UserContact `gorm:"foreignKey:UserID"`
}

// BeforeCreate is a gorm hook that is called before creating a new user
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.ID, err = u.generateID()
	if err != nil {
		return
	}

	return nil
}

// PubKeyObj returns the go-sdk primitives.PublicKey object from the user's PubKey string
func (u *User) PubKeyObj() (*primitives.PublicKey, error) {
	pub, err := primitives.PublicKeyFromString(u.PubKey)
	if err != nil {
		return nil, fmt.Errorf("invalid PubKey: %w", err)
	}
	return pub, nil
}

func (u *User) generateID() (string, error) {
	if u.PubKey == "" {
		return "", fmt.Errorf("PubKey is required")
	}
	pubKey, err := u.PubKeyObj()
	if err != nil {
		return "", err
	}
	addr, err := script.NewAddressFromPublicKey(pubKey, true)
	if err != nil {
		return "", fmt.Errorf("failed to create address from public key: %w", err)
	}
	return addr.AddressString, nil
}
