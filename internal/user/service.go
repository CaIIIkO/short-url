package user

type RepositoryInterface interface {
}

type Service struct {
	repo RepositoryInterface
}

func NewAdService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}
