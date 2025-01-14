package response

import "github.com/bitcoin-sv/spv-wallet/models/bsv"

type UserInfo struct {
	CurrentBalance bsv.Satoshis `json:"currentBalance"`
}
