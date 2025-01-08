package testabilities

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

type mockRepository struct {
	contacts []engine.Contact
}