package mapping

import (
	paymailmodels "github.com/bitcoin-sv/spv-wallet/engine/paymail/models"
	"github.com/bitcoin-sv/spv-wallet/models/response/adminresponse"
)

// CreatedPaymailResponse maps a paymail to a paymail response
func CreatedPaymailResponse(p *paymailmodels.Paymail) adminresponse.Paymail {
	return adminresponse.Paymail{
		ID:         p.ID,
		Alias:      p.Alias,
		Domain:     p.Domain,
		Paymail:    p.Alias + "@" + p.Domain,
		PublicName: p.PublicName,
		Avatar:     p.Avatar,
	}
}
