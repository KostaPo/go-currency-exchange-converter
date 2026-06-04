package currency

import "context"

type Service interface {
	GetAll(ctx context.Context) ([]*Currency, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(ctx context.Context) ([]*Currency, error) {

	currencies, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return currencies, nil
}
