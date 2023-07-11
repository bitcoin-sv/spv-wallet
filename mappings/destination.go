package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux-server/mappings/common"
)

// MapToDestinationContract will map the bux destination model to the bux-models contract
func MapToDestinationContract(d *bux.Destination) *buxmodels.Destination {
	if d == nil {
		return nil
	}

	return &buxmodels.Destination{
		Model:         *common.MapToContract(&d.Model),
		ID:            d.ID,
		XpubID:        d.XpubID,
		LockingScript: d.LockingScript,
		Type:          d.Type,
		Chain:         d.Chain,
		Num:           d.Num,
		Address:       d.Address,
		DraftID:       d.DraftID,
		Monitor:       d.Monitor.Time,
	}
}

// MapToDestinationBux will map the bux-models destination contract to the bux destination model
func MapToDestinationBux(d *buxmodels.Destination) *bux.Destination {
	return &bux.Destination{
		Model:         *common.MapToModel(&d.Model),
		ID:            d.ID,
		XpubID:        d.XpubID,
		LockingScript: d.LockingScript,
		Type:          d.Type,
		Chain:         d.Chain,
		Num:           d.Num,
		Address:       d.Address,
		DraftID:       d.DraftID,
	}
}
