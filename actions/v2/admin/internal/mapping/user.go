package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users/usersmodels"
	"github.com/bitcoin-sv/spv-wallet/models/response/adminresponse"
	"github.com/samber/lo"
)

// CreatedUserResponse maps a user to a response
func UserToResponse(u *usersmodels.User) adminresponse.User {
	return adminresponse.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		PublicKey: u.PublicKey,
		Paymails:  lo.Map(u.Paymails, UsersPaymailToResponse),
	}
}

// UsersPaymailResponse maps a user's paymail to a response
func UsersPaymailToResponse(p *usersmodels.Paymail, _ int) adminresponse.Paymail {
	return adminresponse.Paymail{
		ID:         p.ID,
		Alias:      p.Alias,
		Domain:     p.Domain,
		Paymail:    p.Alias + "@" + p.Domain,
		PublicName: p.PublicName,
		Avatar:     p.Avatar,
	}
}
