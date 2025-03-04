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

// RequestCreateUserToNewUserModel maps a create user request to new user model
func RequestCreateUserToNewUserModel(r *api.RequestsCreateUser) (*usersmodels.NewUser, error) {
	newUser := &usersmodels.NewUser{
		PublicKey: r.PublicKey,
	}

	if isPaymailDefined(r) {
		newPaymail, err := RequestAddPaymailToNewPaymailModel(r.Paymail, "")
		if err != nil {
			return nil, err
		}

		newUser.Paymail = newPaymail
	}

	return newUser, nil
}

func isPaymailDefined(r *api.RequestsCreateUser) bool {
	return r.Paymail != nil
}
