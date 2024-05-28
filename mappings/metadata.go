package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapToSpvWalletMetadata will map the *spvwalletmodels.Metadata to *spv.Metadata
func MapToSpvWalletMetadata(metadata *models.Metadata) *engine.Metadata {
	if metadata == nil {
		return nil
	}

	output := engine.Metadata(*metadata)
	return &output
}

// MapToMetadata converts "raw" key-value map to aliased engine.Metadata
func MapToMetadata(explicitMap *map[string]interface{}) *engine.Metadata {
	if explicitMap == nil {
		return nil
	}
	m := engine.Metadata(*explicitMap)
	return &m
}
