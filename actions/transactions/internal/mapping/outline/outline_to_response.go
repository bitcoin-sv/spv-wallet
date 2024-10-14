package outline

import (
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
	model "github.com/bitcoin-sv/spv-wallet/models/transaction"
	"github.com/mitchellh/mapstructure"
)

// ToResponse converts a draft transaction to a response model.
func ToResponse(tx *draft.Transaction) *model.AnnotatedTransaction {
	res := &model.AnnotatedTransaction{}
	err := mapstructure.Decode(tx, res)
	if err != nil {
		panic(err)
	}
	return res
}
