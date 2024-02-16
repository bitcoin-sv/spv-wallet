package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	spvwalletmodels "github.com/bitcoin-sv/spv-wallet/models"
)

// MapToSPVMetadata will map the *spvwalletmodels.Metadata to *spv.Metadata
func MapToSPVMetadata(metadata *spvwalletmodels.Metadata) *engine.Metadata {
	if metadata == nil {
		return nil
	}

	output := engine.Metadata(*metadata)
	return &output
}
