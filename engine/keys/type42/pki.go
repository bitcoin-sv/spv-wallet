package type42

import (
	"fmt"
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
)

const rotationSuffix = "0"

// PaymailPKI (Public Key Infrastructure) derives a public key using a constant derivation key for provided paymail.
func PaymailPKI(pubKey *primitives.PublicKey, alias, domain string) (*primitives.PublicKey, error) {
	if alias == "" || domain == "" {
		return nil, ErrDeriveKey
	}

	derivationKey := fmt.Sprintf("1-paymail_pki-%s@%s_%s", alias, domain, rotationSuffix)
	derivedPubByRef, err := derive(pubKey, derivationKey)
	if err != nil {
		return nil, ErrDeriveKey.Wrap(err)
	}
	return derivedPubByRef, nil
}
