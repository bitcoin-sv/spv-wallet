package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/engine/domainmodels"
	"github.com/bitcoin-sv/spv-wallet/models/response/adminresponse"
)

// CreatedPaymailResponse maps a paymail to a paymail response
func CreatedPaymailResponse(p *domainmodels.Paymail) adminresponse.Paymail {
	return adminresponse.Paymail{
		ID:         p.ID,
		Alias:      p.Alias,
		Domain:     p.Domain,
		Paymail:    p.Alias + "@" + p.Domain,
		PublicName: p.PublicName,
		Avatar:     p.Avatar,
	}
}
