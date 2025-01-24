package type42

import (
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

var (
	anyonePriv, _ = primitives.PrivateKeyFromBytes([]byte{1})
)

func derive(pubKey *primitives.PublicKey, derivationKey string) (*primitives.PublicKey, error) {
	if pubKey == nil {
		return nil, ErrDeriveKey.Wrap(spverrors.Newf("public key is nil"))
	}
	derivedPubByRef, err := pubKey.DeriveChild(anyonePriv, derivationKey)
	if err != nil {
		return nil, ErrDeriveKey.Wrap(err)
	}
	return derivedPubByRef, nil
}
