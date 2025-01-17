package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/engine/user/usermodels"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/models/response/adminresponse"
)

// CreatedUserResponse maps a user to a user response
func CreatedUserResponse(u *usermodels.User) adminresponse.User {
	return adminresponse.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		PublicKey: u.PublicKey,
		Paymails:  utils.MapSlice(u.Paymails, CreatedPaymailResponse),
	}
}
