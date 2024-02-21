package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux-server/mappings/common"
)

func MapToContactContract(c *bux.Contact) *buxmodels.Contact {
	if c == nil {
		return nil
	}

	return &buxmodels.Contact{
		Model:    *common.MapToContract(&c.Model),
		ID:       c.ID,
		FullName: c.FullName,
		Paymail:  c.Paymail,
		PubKey:   c.PubKey,
		XpubID:   c.XpubID,
		Status:   buxmodels.ContactStatus(c.Status),
	}
}
