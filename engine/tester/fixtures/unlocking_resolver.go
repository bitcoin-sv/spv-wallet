package fixtures

import (
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/custominstructions"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/keys/type42"
)

type UnlockingTemplateResolver struct {
	Template sdk.UnlockingScriptTemplate
}

func (un *UnlockingTemplateResolver) Type42(acc *custominstructions.Accumulator[primitives.PrivateKey], instruction string) (bool, error) {
	priv, err := type42.DerivePrivateKey(acc.Key, instruction)
	if err != nil {
		panic("Invalid setup of user fixture, cannot restore type42 instruction: " + err.Error())
	}
	acc.Key = priv
	return true, nil
}

func (un *UnlockingTemplateResolver) Sign(acc *custominstructions.Accumulator[primitives.PrivateKey], _ string) (bool, error) {
	template, err := p2pkh.Unlock(acc.Key, ptr(sighash.AllForkID))
	if err != nil {
		panic("Invalid setup of user fixture, cannot restore unlocking script: " + err.Error())
	}
	un.Template = template
	return false, nil
}

func (un *UnlockingTemplateResolver) Finalize(acc *custominstructions.Accumulator[primitives.PrivateKey]) error {
	if un.Template == nil {
		// this is implicit "Sign" if there was no "Sign" instruction in provided custom instructions
		_, err := un.Sign(acc, "P2PKH")
		if err != nil {
			return err
		}
	}
	return nil
}
