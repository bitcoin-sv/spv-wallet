package custominstructions

import (
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/custominstructions/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/keys/type42"
)

// NewAddressInterpreter creates a new custom instructions interpreter that resolves a public key to an address.
func NewAddressInterpreter() *Interpreter[*AddressResolver, primitives.PublicKey] {
	return NewInterpreter(&AddressResolver{})
}

// NewLockingScriptInterpreter creates a new custom instructions interpreter that resolves a public key to a locking script (and address).
func NewLockingScriptInterpreter() *Interpreter[*LockingScriptResolver, primitives.PublicKey] {
	return NewInterpreter(&LockingScriptResolver{})
}

// AddressResolver implements resolver for custom instructions that resolve a public key to an address.
type AddressResolver struct {
	Address *script.Address
}

// Type42 derives a new public key from the current public key using a type42 method.
func (ar *AddressResolver) Type42(acc *Accumulator[primitives.PublicKey], instruction string) (bool, error) {
	pub, err := type42.Derive(acc.Key, instruction)
	if err != nil {
		return false, errors.ErrType42DerivationFailed.Wrap(err)
	}
	acc.Key = pub
	return true, nil
}

// Sign derives an address from the current public key.
func (ar *AddressResolver) Sign(acc *Accumulator[primitives.PublicKey], _ string) (bool, error) {
	addr, err := script.NewAddressFromPublicKey(acc.Key, true)
	if err != nil {
		return false, errors.ErrGettingAddressFromPublicKey.Wrap(err)
	}
	ar.Address = addr
	return false, nil
}

// Finalize ensures that the address is resolved.
func (ar *AddressResolver) Finalize(acc *Accumulator[primitives.PublicKey]) error {
	if ar.Address == nil {
		// this is implicit "Sign" if there was no "Sign" instruction in provided custom instructions
		_, err := ar.Sign(acc, "P2PKH")
		if err != nil {
			return spverrors.Wrapf(err, "failed to finalize and derive address")
		}
	}
	return nil
}

// LockingScriptResolver implements resolver for custom instructions that resolve a public key to a locking script (and address).
// Uses the AddressResolver because a locking script is easily derived from an address.
type LockingScriptResolver struct {
	AddressResolver
	LockingScript *script.Script
}

// Finalize ensures that the address is resolved and derives the locking script from the address.
func (lr *LockingScriptResolver) Finalize(acc *Accumulator[primitives.PublicKey]) error {
	err := lr.AddressResolver.Finalize(acc)
	if err != nil {
		return err
	}

	lockingScript, err := p2pkh.Lock(lr.Address)
	if err != nil {
		return errors.ErrGettingLockingScript.Wrap(err)
	}
	lr.LockingScript = lockingScript

	return nil
}
