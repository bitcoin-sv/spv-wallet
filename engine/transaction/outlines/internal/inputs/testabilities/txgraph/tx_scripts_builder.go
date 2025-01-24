package txgraph

import (
	"testing"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/stretchr/testify/require"
)

type TxScriptsBuilder struct {
	t          *testing.T
	paymail    string
	privateKey *primitives.PrivateKey
	publicKey  *primitives.PublicKey
}

func (tx *TxScriptsBuilder) PrivateKey() *primitives.PrivateKey { return tx.privateKey }

func (tx *TxScriptsBuilder) PublicKey() *primitives.PublicKey { return tx.publicKey }

func (tx *TxScriptsBuilder) P2PKHUnlockingScriptTemplate() *p2pkh.P2PKH {
	tx.t.Helper()

	script, err := p2pkh.Unlock(tx.privateKey, nil)
	require.NoError(tx.t, err, "failed to return unlocking script")
	return script
}

func (tx *TxScriptsBuilder) P2PKHLockingScript() *script.Script {
	tx.t.Helper()

	addr, err := script.NewAddressFromPublicKey(tx.publicKey, true)
	require.NoError(tx.t, err, "failed to return addr from pub key")

	script, err := p2pkh.Lock(addr)
	require.NoError(tx.t, err, "failed to return locking script")
	return script
}

func NewTxScriptsBuilder(t *testing.T, xPriv, paymail string) *TxScriptsBuilder {
	t.Helper()

	hdKey, err := bip32.GenerateHDKeyFromString(xPriv)
	require.NoErrorf(t, err, "failed to generate HD Key from string %s", xPriv)

	privateKey, err := bip32.GetPrivateKeyFromHDKey(hdKey)
	require.NoError(t, err, "failed to retrive priv key from HD key")

	return &TxScriptsBuilder{
		t:          t,
		paymail:    paymail,
		privateKey: privateKey,
		publicKey:  privateKey.PubKey(),
	}
}
