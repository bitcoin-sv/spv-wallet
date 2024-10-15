package outline

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
	model "github.com/bitcoin-sv/spv-wallet/models/transaction"
	"github.com/mitchellh/mapstructure"
)

// ToResponse converts a draft transaction to a response model.
func ToResponse(tx *outlines.Transaction) (*model.AnnotatedTransaction, error) {
	res := &model.AnnotatedTransaction{}
	err := mapstructure.Decode(tx, res)
	if err != nil {
		return nil, spverrors.ErrCannotMapFromModel.Wrap(err)
	}
	return res, nil
}
