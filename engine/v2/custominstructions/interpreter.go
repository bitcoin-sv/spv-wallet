package custominstructions

import (
	"github.com/bitcoin-sv/spv-wallet/engine/v2/custominstructions/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

type Interpreter[R Resolver[TKey], TKey InputKeys] struct {
	resolver R
}

func NewInterpreter[R Resolver[TKey], TKey InputKeys](resolver R) *Interpreter[R, TKey] {
	return &Interpreter[R, TKey]{
		resolver: resolver,
	}
}

func (p *Interpreter[I, TKey]) Process(publicKey *TKey, instructions bsv.CustomInstructions) (I, error) {
	acc := &Accumulator[TKey]{
		Key: publicKey,
	}
	var err error
	var proceed bool
	for _, instruction := range instructions {
		switch instruction.Type {
		case Type42:
			proceed, err = p.resolver.Type42(acc, instruction.Instruction)
		case Sign:
			proceed, err = p.resolver.Sign(acc, instruction.Instruction)
		default:
			return p.resolver, errors.ErrUnknownInstructionType
		}
		if err != nil {
			return p.resolver, errors.ErrProcessingCustomInstructions.Wrap(err)
		}
		if !proceed {
			break
		}
	}

	err = p.resolver.Finalize(acc)
	if err != nil {
		return p.resolver, errors.ErrFinalizingCustomInstructions.Wrap(err)
	}

	return p.resolver, nil
}
