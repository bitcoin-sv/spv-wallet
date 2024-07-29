// Package mappings is a package that contains the mappings for the access keys package.
package mappings

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// MapToOldAccessKeyContract will map the access key to the spv-wallet-models contract
func MapToOldAccessKeyContract(ac *engine.AccessKey) *models.AccessKey {
	if ac == nil {
		return nil
	}

	var revokedAt *time.Time
	if !ac.RevokedAt.IsZero() {
		revokedAt = &ac.RevokedAt.Time
	}

	return &models.AccessKey{
		Model:     *common.MapToContract(&ac.Model),
		ID:        ac.ID,
		XpubID:    ac.XpubID,
		RevokedAt: revokedAt,
		Key:       ac.Key,
	}
}

// MapToAccessKeyContract will map the access key to the spv-wallet-models contract
func MapToAccessKeyContract(ac *engine.AccessKey) *response.AccessKey {
	if ac == nil {
		return nil
	}

	var revokedAt *time.Time
	if !ac.RevokedAt.IsZero() {
		revokedAt = &ac.RevokedAt.Time
	}

	return &response.AccessKey{
		Model:     *common.MapToContract(&ac.Model),
		ID:        ac.ID,
		XpubID:    ac.XpubID,
		RevokedAt: revokedAt,
		Key:       ac.Key,
	}
}
