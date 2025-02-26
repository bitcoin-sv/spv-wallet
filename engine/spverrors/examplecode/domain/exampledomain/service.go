package exampledomain

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	domainerr "github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/domain/exampledomain/errors"
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

func (s *Service) DoSth(fail *api.ModelsFailingPoint) ([]string, error) {
	if fail != nil {
		switch *fail {
		case api.DomainIllegalArgument:
			// NOTE: at domain level we don't know if the argument was provided by external client
			// or by other services which uses this service.
			return nil, domainerr.WrongArgument.
				New("Argument has from format %s", "abc")
		case api.ArcDoubleSpentAttempt:
			return nil, domainerr.SomeARCError.
				New("Double spent attempt").
				WithProperty(errdef.PropPublicHint, "This is a public hint which can help debug what ARC says")
		}
	}

	arr, err := s.repo.Search(fail)
	if err != nil {
		return nil, errorx.Decorate(err, "Cannot fetch data from database")
	}
	return arr, nil
}
