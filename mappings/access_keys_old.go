// Package mappings is a package that contains the mappings for the access keys package.
package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapToOldAccessKeyContract will map the access key to the spv-wallet-models contract
func MapToOldAccessKeyContract(ac *engine.AccessKey) *models.AccessKey {
	if ac == nil {
		return nil
	}

	return &models.AccessKey{
		Model:  *common.MapToOldContract(&ac.Model),
		ID:     ac.ID,
		XpubID: ac.XpubID,
		Key:    ac.Key,
	}
}
