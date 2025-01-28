package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users/usersmodels"
	"github.com/bitcoin-sv/spv-wallet/mapper"
	"github.com/bitcoin-sv/spv-wallet/models/response/adminresponse"
)

// CreatedUserResponse maps a user to a response
func UserToResponse(u *usersmodels.User) adminresponse.User {
	return adminresponse.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		PublicKey: u.PublicKey,
		Paymails:  mapper.MapSlice(u.Paymails, UsersPaymailToResponse),
	}
}

// UsersPaymailResponse maps a user's paymail to a response
func UsersPaymailToResponse(p *usersmodels.Paymail) adminresponse.Paymail {
	return adminresponse.Paymail{
		ID:         p.ID,
		Alias:      p.Alias,
		Domain:     p.Domain,
		Paymail:    p.Alias + "@" + p.Domain,
		PublicName: p.PublicName,
		Avatar:     p.Avatar,
	}
}
