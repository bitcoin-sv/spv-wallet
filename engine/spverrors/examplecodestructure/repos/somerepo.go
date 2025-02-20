//go:build errorx
// +build errorx

package repos

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/spike/domain/domainerr"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/spike/domain/transaction"
)

type FailingRepo struct {
}

func (f *FailingRepo) Search(query string) ([]*transaction.Transaction, error) {
	if query == "success" {
		return []*transaction.Transaction{
			{
				ID: "1",
			},
		}, nil
	}
	return nil, domainerr.QueryFailed.New("search transactions failed")
}

func (f *FailingRepo) Find(id string) (*transaction.Transaction, error) {
	return &transaction.Transaction{
		ID: id,
	}, nil
}

func (f *FailingRepo) Save(t *transaction.Transaction) error {
	return domainerr.SaveFailed.New("save transaction failed")
}

func NewRepo() *FailingRepo {
	return &FailingRepo{}
}
