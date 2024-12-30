package database

import (
	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	dberrors "github.com/bitcoin-sv/spv-wallet/engine/database/errors"
	"time"
)

type User struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Xpub string //TODO: change to Pub
}

// UserPub -> PKI ->

func (u *User) GetHDPublicKey() (*bip32.ExtendedKey, error) {
	hdPub, err := bip32.NewKeyFromString(u.Xpub)
	if err != nil {
		return nil, dberrors.ErrConvertTOHDPubKey.Wrap(err)
	}
	return hdPub, nil
}

func (u *User) GetPublicKey() (*primitives.PublicKey, error) {
	hdPub, err := u.GetHDPublicKey()
	pub, err := bip32.GetPublicKeyFromHDKey(hdPub)
	if err != nil {
		return nil, dberrors.ErrConvertToPubKey.Wrap(err)
	}
	return pub, nil
}
