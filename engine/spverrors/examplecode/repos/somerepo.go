package repos

import (
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/api"
	repoerr "github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/repos/errors"
)

type FailingRepo struct{}

func (f *FailingRepo) Search(fail *api.ModelsFailingPoint) ([]string, error) {
	if fail == nil {
		return []string{"success"}, nil
	}

	switch *fail {
	case api.DbConnection:
		err := fmt.Errorf("db connection failed from external lib")
		return nil, repoerr.DbConnectionFailed.Wrap(err, "db connection failed")
	case api.DbQuery:
		err := fmt.Errorf("query failed from external lib")
		return nil, repoerr.DbQueryFailed.Wrap(err, "query failed")
	default:
		return nil, repoerr.DbShouldNeverHappen.NewWithNoMessage()
	}
}

func NewRepo() *FailingRepo {
	return &FailingRepo{}
}
