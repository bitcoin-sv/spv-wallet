package destinations

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

type CreateDestination struct {
	Metadata engine.Metadata `json:"metadata"`
}

type UpdateDestination struct {
	Id            string          `json:"id"`
	Address       string          `json:"address"`
	LockingScript string          `json:"locking_script"`
	Metadata      engine.Metadata `json:"metadata"`
}
