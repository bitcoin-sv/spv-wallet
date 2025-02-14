package type42

import (
	"fmt"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// Destination represents a destination for a transaction output.
type Destination struct {
	PubKey        *primitives.PublicKey
	ReferenceID   string
	DerivationKey string
}

// NewDestinationWithReference derives a public key using a reference ID.
// It is intended to be used to derive a public key for paymail destinations.
func NewDestinationWithReference(pubKey *primitives.PublicKey, referenceID string) (Destination, error) {
	dst := Destination{
		ReferenceID: referenceID,
	}
	var err error
	if dst.ReferenceID == "" {
		return dst, ErrDeriveKey.Wrap(spverrors.Newf("reference ID is empty"))
	}

	dst.DerivationKey = fmt.Sprintf("1-destination-%s", referenceID)
	dst.PubKey, err = derive(pubKey, dst.DerivationKey)
	if err != nil {
		return dst, err
	}
	return dst, nil
}

// NewDestinationWithRandomReference helps to generate a destination with a random reference ID.
func NewDestinationWithRandomReference(pubKey *primitives.PublicKey) (Destination, error) {
	referenceID, err := utils.RandomHex(16)
	if err != nil {
		return Destination{}, ErrRandomReferenceID.Wrap(err)
	}

	return NewDestinationWithReference(pubKey, referenceID)
}
