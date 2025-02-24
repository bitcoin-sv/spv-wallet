package custominstructions

import (
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/custominstructions/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/keys/type42"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

func Address(publicKey primitives.PublicKey, instructions bsv.CustomInstructions) (*script.Address, error) {
	pub := &publicKey
	var err error
	for _, instruction := range instructions {
		switch instruction.Type {
		case Type42:
			pub, err = type42.Derive(pub, instruction.Instruction)
			if err != nil {
				return nil, errors.ErrType42DerivationFailed.Wrap(err)
			}
		default:
			return nil, errors.ErrUnknownInstructionType
		}
	}

	addr, err := script.NewAddressFromPublicKey(pub, true)
	if err != nil {
		return nil, errors.ErrGettingAddressFromPublicKey.Wrap(err)
	}

	return addr, nil
}
