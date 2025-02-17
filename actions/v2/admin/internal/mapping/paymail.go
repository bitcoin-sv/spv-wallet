package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
)

// PaymailToAdminResponse maps a paymail to a response
func PaymailToAdminResponse(p *paymailsmodels.Paymail) api.ModelsPaymail {
	return api.ModelsPaymail{
		Id:         p.ID,
		Alias:      p.Alias,
		Domain:     p.Domain,
		Paymail:    p.Alias + "@" + p.Domain,
		PublicName: p.PublicName,
		Avatar:     p.Avatar,
	}
}
