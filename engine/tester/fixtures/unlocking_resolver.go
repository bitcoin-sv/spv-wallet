package fixtures

import (
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/keys/type42"
	"github.com/samber/lo"
)

type UnlockingTemplateResolver struct {
	Template sdk.UnlockingScriptTemplate
	privKey  *primitives.PrivateKey
}

func (un *UnlockingTemplateResolver) Initialize(privKey *primitives.PrivateKey) error {
	un.privKey = privKey
	return nil
}

func (un *UnlockingTemplateResolver) Type42(instruction string) (bool, error) {
	priv, err := type42.DerivePrivateKey(un.privKey, instruction)
	if err != nil {
		panic("Invalid setup of user fixture, cannot restore type42 instruction: " + err.Error())
	}
	un.privKey = priv
	return true, nil
}

func (un *UnlockingTemplateResolver) Sign(_ string) (bool, error) {
	template, err := p2pkh.Unlock(un.privKey, lo.ToPtr(sighash.AllForkID))
	if err != nil {
		panic("Invalid setup of user fixture, cannot restore unlocking script: " + err.Error())
	}
	un.Template = template
	return false, nil
}

func (un *UnlockingTemplateResolver) Finalize() error {
	if un.Template == nil {
		// this is implicit "Sign" if there was no "Sign" instruction in provided custom instructions
		_, _ = un.Sign("P2PKH")
	}
	return nil
}
