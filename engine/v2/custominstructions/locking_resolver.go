package custominstructions

import (
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/custominstructions/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/keys/type42"
)

func NewLockingScriptInterpreter() *Interpreter[*LockingScriptResolver, primitives.PublicKey] {
	return NewInterpreter(&LockingScriptResolver{})
}

type LockingScriptResolver struct {
	Address       *script.Address
	LockingScript *script.Script
}

func (ar *LockingScriptResolver) Type42(acc *Accumulator[primitives.PublicKey], instruction string) (bool, error) {
	pub, err := type42.Derive(acc.Key, instruction)
	if err != nil {
		return false, errors.ErrType42DerivationFailed.Wrap(err)
	}
	acc.Key = pub
	return true, nil
}

func (ar *LockingScriptResolver) Sign(acc *Accumulator[primitives.PublicKey], _ string) (bool, error) {
	addr, err := script.NewAddressFromPublicKey(acc.Key, true)
	if err != nil {
		return false, errors.ErrGettingAddressFromPublicKey.Wrap(err)
	}
	ar.Address = addr

	lockingScript, err := p2pkh.Lock(addr)
	if err != nil {
		return false, errors.ErrGettingLockingScript.Wrap(err)
	}
	ar.LockingScript = lockingScript

	return false, nil
}

func (ar *LockingScriptResolver) Finalize(acc *Accumulator[primitives.PublicKey]) error {
	if ar.Address == nil {
		// this is implicit "Sign" if there was no "Sign" instruction in provided custom instructions
		_, err := ar.Sign(acc, "P2PKH")
		if err != nil {
			return err
		}
	}
	return nil
}
