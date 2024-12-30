package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/models/response/adminresponse"
)

// UserResponse maps a user to a user response
func UserResponse(u *database.User) adminresponse.User {
	return adminresponse.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		PublicKey: u.PubKey,
		Paymails:  mapSlice(PaymailResponse, u.Paymails),
	}
}
