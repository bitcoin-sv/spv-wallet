package accesskeys

import "github.com/bitcoin-sv/spv-wallet/engine"

type CreateAccessKey struct {
	Metadata engine.Metadata `json:"metadata"`
}
