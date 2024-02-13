package mappings

import (
	"github.com/BuxOrg/bux"
	spvwalletmodels "github.com/BuxOrg/bux-models"
)

// MapToSPVMetadata will map the *spvwalletmodels.Metadata to *spv.Metadata
func MapToSPVMetadata(metadata *spvwalletmodels.Metadata) *bux.Metadata {
	if metadata == nil {
		return nil
	}

	output := bux.Metadata(*metadata)
	return &output
}
