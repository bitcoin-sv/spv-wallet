package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapOldToDestinationContract will map the spv-wallet destination model to the spv-wallet-models contract
func MapOldToDestinationContract(d *engine.Destination) *models.Destination {
	if d == nil {
		return nil
	}

	return &models.Destination{
		Model:                        *common.MapToOldContract(&d.Model),
		ID:                           d.ID,
		XpubID:                       d.XpubID,
		LockingScript:                d.LockingScript,
		Type:                         d.Type,
		Chain:                        d.Chain,
		Num:                          d.Num,
		PaymailExternalDerivationNum: d.PaymailExternalDerivationNum,
		Address:                      d.Address,
		DraftID:                      d.DraftID,
	}
}

// MapOldDestinationModelToEngine will map the spv-wallet-models destination contract to the spv-wallet destination model
func MapOldDestinationModelToEngine(d *models.Destination) *engine.Destination {
	if d == nil {
		return nil
	}

	return &engine.Destination{
		Model:                        *common.MapOldContractToModel(&d.Model),
		ID:                           d.ID,
		XpubID:                       d.XpubID,
		LockingScript:                d.LockingScript,
		Type:                         d.Type,
		Chain:                        d.Chain,
		Num:                          d.Num,
		PaymailExternalDerivationNum: d.PaymailExternalDerivationNum,
		Address:                      d.Address,
		DraftID:                      d.DraftID,
	}
}
