// Package pike provides functionality to work with Pay-to-PubKey-Hash (P2PKH) scripts in a blockchain context.
//
// P2PKH is a common type of Bitcoin transaction that locks funds to a specific public key hash, requiring the
// corresponding private key to produce a valid signature for spending the funds. This package offers utilities
// to create, decode, and validate P2PKH scripts.
//
// The package includes:
// - Functions to generate P2PKH addresses from public keys.
// - Methods to construct P2PKH scriptPubKey and scriptSig scripts.
// - Utilities to decode and inspect P2PKH scripts and addresses.
// - Validation functions to ensure the integrity and correctness of P2PKH scripts.
//
// This package is intended for developers working with Bitcoin or other cryptocurrencies that support P2PKH transactions.
// It abstracts the low-level details and provides a high-level interface for creating and handling P2PKH scripts.
package pike

import (
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/engine/script/template"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/types/type42"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bt/v2/bscript"
)

// GenerateOutputsTemplate creates a Pike output template
func GenerateOutputsTemplate(satoshis uint64) ([]*template.OutputTemplate, error) {
	p2pkhTemplate, err := template.P2PKH(satoshis)
	if err != nil {
		return nil, spverrors.Wrapf(err, "error creating P2PKH template")
	}
	return []*template.OutputTemplate{p2pkhTemplate}, nil
}

// GenerateLockingScriptsFromTemplates converts Pike outputs templates to scripts
func GenerateLockingScriptsFromTemplates(outputsTemplate []*template.OutputTemplate, senderPubKey, receiverPubKey *bec.PublicKey, reference string) ([]string, error) {
	lockingScripts := make([]string, len(outputsTemplate))

	for idx, output := range outputsTemplate {
		templateScript, err := bscript.NewFromHexString(output.Script)
		if err != nil {
			return nil, spverrors.Wrapf(err, "error creating script from hex string")
		}

		dPK, err := type42.DeriveLinkedKey(senderPubKey, receiverPubKey, fmt.Sprintf("%s-%d", reference, idx))
		if err != nil {
			return nil, spverrors.Wrapf(err, "error deriving linked key")
		}

		scriptBytes, err := template.Evaluate(*templateScript, dPK)
		if err != nil {
			return nil, spverrors.Wrapf(err, "error evaluating template script")
		}

		finalScript := bscript.Script(scriptBytes)

		lockingScripts[idx] = finalScript.String()
	}

	return lockingScripts, nil
}
