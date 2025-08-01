package user

type RepositoryInterface interface {
}

type Service struct {
	repo RepositoryInterface
}

func NewUserService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}
