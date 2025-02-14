package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users/usersmodels"
	"github.com/bitcoin-sv/spv-wallet/lox"
	"github.com/samber/lo"
)

// UserToResponse maps a user to a response
func UserToResponse(u *usersmodels.User) api.ModelsUser {
	return api.ModelsUser{
		Id:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		PublicKey: u.PublicKey,
		Paymails:  lo.Map(u.Paymails, lox.MappingFn(UsersPaymailToResponse)),
	}
}

// UsersPaymailToResponse maps a user's paymail to a response
func UsersPaymailToResponse(p *usersmodels.Paymail) api.ModelsPaymail {
	return api.ModelsPaymail{
		Id:         p.ID,
		Alias:      p.Alias,
		Domain:     p.Domain,
		Paymail:    p.Alias + "@" + p.Domain,
		PublicName: p.PublicName,
		Avatar:     p.Avatar,
	}
}
