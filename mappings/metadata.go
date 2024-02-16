package mappings

import (
	"github.com/bitcoin-sv/bux"
	spvwalletmodels "github.com/bitcoin-sv/bux-models"
)

// MapToSPVMetadata will map the *spvwalletmodels.Metadata to *spv.Metadata
func MapToSPVMetadata(metadata *spvwalletmodels.Metadata) *bux.Metadata {
	if metadata == nil {
		return nil
	}

	output := bux.Metadata(*metadata)
	return &output
}
