package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

// MapToMetadata converts "raw" key-value map to aliased engine.Metadata
func MapToMetadata(explicitMap map[string]interface{}) *engine.Metadata {
	if explicitMap == nil {
		return nil
	}
	m := engine.Metadata(explicitMap)
	return &m
}
