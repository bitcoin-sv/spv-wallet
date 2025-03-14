package testabilities

import (
	"testing"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/stretchr/testify/require"
)

// txScripts is a utility structure that provides access to transaction-related scripts
// and keys for creating and validating Bitcoin transactions.
type txScripts struct {
	t          *testing.T             // Test context used for error handling.
	privateKey *primitives.PrivateKey // Private key for signing transactions.
	publicKey  *primitives.PublicKey  // Public key derived from the private key.
}

// p2pKHUnlockingScriptTemplate generates a P2PKH unlocking script template using the private key.
// This template can be used to sign transactions.
func (tx *txScripts) p2pKHUnlockingScriptTemplate() *p2pkh.P2PKH {
	tx.t.Helper()

	unlocking, err := p2pkh.Unlock(tx.privateKey, nil)
	require.NoError(tx.t, err, "Failed to generate unlocking script")
	return unlocking
}

// p2pKHLockingScript generates a P2PKH locking script using the associated public key.
// This script is used to lock outputs to the public key.
func (tx *txScripts) p2pKHLockingScript() *script.Script {
	tx.t.Helper()

	addr, err := script.NewAddressFromPublicKey(tx.publicKey, true)
	require.NoError(tx.t, err, "Failed to generate address from public key")

	locking, err := p2pkh.Lock(addr)
	require.NoError(tx.t, err, "Failed to generate locking script")
	return locking
}

// newTxScripts initializes a new txScripts instance using the provided extended private key (xPriv) and paymail.
// It derives the private and public keys from the xPriv and prepares the structure for generating scripts.
func newTxScripts(t *testing.T, xPriv string) *txScripts {
	t.Helper()

	// Generate the hierarchical deterministic (HD) key from the xPriv string.
	hdKey, err := bip32.GenerateHDKeyFromString(xPriv)
	require.NoErrorf(t, err, "Failed to generate HD key from string: %s", xPriv)

	// Retrieve the private key from the HD key.
	privateKey, err := bip32.GetPrivateKeyFromHDKey(hdKey)
	require.NoError(t, err, "Failed to retrieve private key from HD key")

	// Return a new TxScripts instance.
	return &txScripts{
		t:          t,
		privateKey: privateKey,
		publicKey:  privateKey.PubKey(),
	}
}
