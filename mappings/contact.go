package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func MapToContactContract(c *engine.Contact) *models.Contact {
	if c == nil {
		return nil
	}

	return &models.Contact{
		Model:    *common.MapToContract(&c.Model),
		ID:       c.ID,
		FullName: c.FullName,
		Paymail:  c.Paymail,
		PubKey:   c.PubKey,
		XpubID:   c.XpubID,
		Status:   models.ContactStatus(c.Status),
	}
}
