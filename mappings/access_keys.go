// Package mappings is a package that contains the mappings for the access keys package.
package mappings

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapToAccessKeyContract will map the access key to the spv-wallet-models contract
func MapToAccessKeyContract(ac *engine.AccessKey) *models.AccessKey {
	if ac == nil {
		return nil
	}

	var revokedAt *time.Time
	if !ac.RevokedAt.IsZero() {
		revokedAt = &ac.RevokedAt.Time
	}

	return &models.AccessKey{
		Model:     *common.MapToOldContract(&ac.Model),
		ID:        ac.ID,
		XpubID:    ac.XpubID,
		RevokedAt: revokedAt,
		Key:       ac.Key,
	}
}
