//go:build errorx
// +build errorx

package transaction

import "github.com/joomcode/errorx"

var UnexpectedQueryValue = errorx.IllegalArgument.NewSubtype("unexpected_query_value")

type Transaction struct {
	ID string
}

type Repo interface {
	Search(string) ([]*Transaction, error)
	Find(id string) (*Transaction, error)
	Save(t *Transaction) error
}

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Search(query string) ([]*Transaction, error) {
	if query == "invalid" {
		return nil, UnexpectedQueryValue.New("query cannot be 'fail'")
	}

	result, err := s.repo.Search(query)
	if err != nil {
		return nil, errorx.Decorate(err, "searching for transactions failed")
	}

	return result, nil
}
