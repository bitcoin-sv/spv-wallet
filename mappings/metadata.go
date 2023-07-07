package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
)

// MapToBuxMetadata will map the *buxmodels.Metadata to *bux.Metadata
func MapToBuxMetadata(metadata *buxmodels.Metadata) *bux.Metadata {
	if metadata == nil {
		return nil
	}

	output := bux.Metadata(*metadata)
	return &output
}
