package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux-server/mappings/common"
)

func MapToDestinationContract(d *bux.Destination) *buxmodels.Destination {
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
