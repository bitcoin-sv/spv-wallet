package exampledomain

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/errdef"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/repos"
	"github.com/joomcode/errorx"
)

type Transaction struct {
	ID string
}

type Service struct {
	repo *repos.FailingRepo
}

func NewService(repo *repos.FailingRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Search(fail *api.ModelsFailingPoint) ([]string, error) {
	arr, err := s.repo.Search(fail)
	if err != nil {
		if errorx.HasTrait(err, errdef.TraitIllegalArgument) {
			return nil, errdef.AsClientError(err, errdef.ClientUnprocessableEntity)
		}
		return nil, errorx.Decorate(err, "Cannot fetch data from database")
	}
	return arr, nil
}
