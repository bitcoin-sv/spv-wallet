package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/engine/domainmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/models/response/adminresponse"
)

// CreatedUserResponse maps a user to a user response
func CreatedUserResponse(u *domainmodels.User) adminresponse.User {
	return adminresponse.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		PublicKey: u.PublicKey,
		Paymails:  utils.MapSlice(u.Paymails, CreatedPaymailResponse),
	}
}
