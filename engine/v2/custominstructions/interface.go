package custominstructions

import (
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
)

type InputKeys interface {
	primitives.PublicKey | primitives.PrivateKey
}

type Accumulator[TKey InputKeys] struct {
	Key *TKey
}

type Resolver[TKey InputKeys] interface {
	Type42(acc *Accumulator[TKey], instruction string) (proceed bool, err error)
	Sign(acc *Accumulator[TKey], instruction string) (proceed bool, err error)

	Finalize(acc *Accumulator[TKey]) error
}
