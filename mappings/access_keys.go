// Package mappings is a package that contains the mappings for the access keys package.
package mappings

import (
	"time"

	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux-server/mappings/common"
)

// MapToAccessKeyContract will map the access key to the bux-models contract
func MapToAccessKeyContract(ac *bux.AccessKey) *buxmodels.AccessKey {
	if ac == nil {
		return nil
	}

	var revokedAt *time.Time
	if !ac.RevokedAt.IsZero() {
		revokedAt = &ac.RevokedAt.Time
	}

	return &buxmodels.AccessKey{
		Model:     *common.MapToContract(&ac.Model),
		ID:        ac.ID,
		XpubID:    ac.XpubID,
		RevokedAt: revokedAt,
		Key:       ac.Key,
	}
}
