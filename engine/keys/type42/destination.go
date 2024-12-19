package type42

import (
	"fmt"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// Destination derives a public key using a reference ID.
// It is intended to be used to derive a public key for paymail destinations.
func Destination(pubKey *primitives.PublicKey, referenceID string) (*primitives.PublicKey, error) {
	if referenceID == "" {
		return nil, ErrDeriveKey.Wrap(spverrors.Newf("reference ID is empty"))
	}
	derivedPubByRef, err := derive(pubKey, fmt.Sprintf("1-destination-%s", referenceID))
	if err != nil {
		return nil, err
	}
	return derivedPubByRef, nil
}
