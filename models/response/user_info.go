package response

import "github.com/bitcoin-sv/spv-wallet/models/bsv"

// UserInfo represents the response model for current user information
type UserInfo struct {
	CurrentBalance bsv.Satoshis `json:"currentBalance"`
}
