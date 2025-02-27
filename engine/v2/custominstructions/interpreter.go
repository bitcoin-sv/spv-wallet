package custominstructions

import (
	"github.com/bitcoin-sv/spv-wallet/engine/v2/custominstructions/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// Interpreter is a struct that is used to interpret custom instructions.
type Interpreter[R Resolver[TKey], TKey InputKeys] struct {
	resolver R
}

// NewInterpreter creates a new interpreter for custom instructions with a given resolver.
func NewInterpreter[R Resolver[TKey], TKey InputKeys](resolver R) *Interpreter[R, TKey] {
	return &Interpreter[R, TKey]{
		resolver: resolver,
	}
}

// Process processes custom instructions for a given key.
func (p *Interpreter[I, TKey]) Process(key *TKey, instructions bsv.CustomInstructions) (I, error) {
	var err error
	var proceed bool

	err = p.resolver.Initialize(key)
	if err != nil {
		return p.resolver, errors.ErrInitializingCustomInstructions.Wrap(err)
	}

	for _, instruction := range instructions {
		switch instruction.Type {
		case Type42:
			proceed, err = p.resolver.Type42(instruction.Instruction)
		case Sign:
			proceed, err = p.resolver.Sign(instruction.Instruction)
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

	err = p.resolver.Finalize()
	if err != nil {
		return p.resolver, errors.ErrFinalizingCustomInstructions.Wrap(err)
	}

	return p.resolver, nil
}
