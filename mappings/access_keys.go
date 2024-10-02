// Package mappings is a package that contains the mappings for the access keys package.
package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// MapToAccessKeyContract will map the access key to the spv-wallet-models contract
func MapToAccessKeyContract(ac *engine.AccessKey) *response.AccessKey {
	if ac == nil {
		return nil
	}

	return &response.AccessKey{
		Model:  *common.MapToContract(&ac.Model),
		ID:     ac.ID,
		XpubID: ac.XpubID,
		Key:    ac.Key,
	}
}
