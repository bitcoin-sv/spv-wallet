package type42

import primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"

const pkiDerivationKey = "1-pki-0"

// PKI (Public Key Infrastructure) derives a public key using a constant derivation key.
func PKI(pubKey *primitives.PublicKey) (*primitives.PublicKey, error) {
	derivedPubByRef, err := derive(pubKey, pkiDerivationKey)
	if err != nil {
		return nil, ErrDeriveKey.Wrap(err)
	}
	return derivedPubByRef, nil
}
