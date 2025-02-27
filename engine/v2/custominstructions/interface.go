package custominstructions

import (
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
)

// InputKeys is a union type that can be either a public or private key.
// It is used to define the type of the key that is being used as input for the interpreter.
type InputKeys interface {
	primitives.PublicKey | primitives.PrivateKey
}

// Resolver is an interface that is used to resolve custom instructions.
// It is used by the interpreter.
type Resolver[TKey InputKeys] interface {
	Type42(instruction string) (proceed bool, err error)
	Sign(instruction string) (proceed bool, err error)

	Initialize(key *TKey) error
	Finalize() error
}
